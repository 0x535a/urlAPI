"""Publish archives, SHA256 checksums, and release notes using CI_JOB_TOKEN."""

import hashlib
import json
import os
from pathlib import Path
import subprocess
from urllib.error import HTTPError
from urllib.parse import quote
from urllib.request import Request, urlopen


def main():
    tag = os.environ["CI_COMMIT_TAG"]
    api = f"{os.environ['CI_API_V4_URL']}/projects/{os.environ['CI_PROJECT_ID']}"
    headers = {"JOB-TOKEN": os.environ["CI_JOB_TOKEN"]}

    def request(method, url, data=None, content_type="application/json"):
        req = Request(url, data=data, method=method,
                      headers={**headers, "Content-Type": content_type})
        with urlopen(req, timeout=300) as response:
            return response.read()

    archives = sorted(
        Path("packages") / f"urlAPI_{tag}_{goos}_{arch}.{suffix}"
        for goos in ("linux", "windows")
        for arch in ("amd64", "arm64")
        for suffix in ("zip" if goos == "windows" else "tar.gz",)
    )
    missing = [str(path) for path in archives if not path.is_file()]
    if missing:
        raise SystemExit(f"Missing release archives: {', '.join(missing)}")
    checksums = Path("packages/checksums.txt")
    checksums.write_text("".join(
        f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.name}\n"
        for path in archives
    ))

    previous = subprocess.run(
        ["git", "describe", "--tags", "--match", "v*", "--abbrev=0", f"{tag}^"],
        capture_output=True, text=True,
    )
    revision = f"{previous.stdout.strip()}..{tag}" if previous.returncode == 0 else tag
    notes = subprocess.check_output(
        ["git", "log", "--format=- %s (%h)", revision], text=True,
    ).strip()
    links = []
    for path in [*archives, checksums]:
        url = f"{api}/packages/generic/urlapi/{quote(tag, safe='')}/{path.name}"
        request("PUT", url, path.read_bytes(), "application/octet-stream")
        links.append({"name": path.name, "url": url, "link_type": "package"})
        print(f"Uploaded {path.name}")

    release_url = f"{api}/releases/{quote(tag, safe='')}"
    payload = {"name": tag, "description": f"urlAPI {tag}\n\n{notes}",
               "assets": {"links": links}}
    try:
        existing_release = json.loads(request("GET", release_url))
    except HTTPError as error:
        if error.code != 404:
            raise
        payload["tag_name"] = tag
        request("POST", f"{api}/releases", json.dumps(payload).encode())
    else:
        payload.pop("assets")
        request("PUT", release_url, json.dumps(payload).encode())
        existing_links = {
            link["name"]: link for link in existing_release.get("assets", {}).get("links", [])
        }
        for link in links:
            existing = existing_links.get(link["name"])
            if existing:
                if existing.get("url") != link["url"] or existing.get("link_type") != "package":
                    request("PUT", f"{release_url}/assets/links/{existing['id']}",
                            json.dumps(link).encode())
            else:
                request("POST", f"{release_url}/assets/links", json.dumps(link).encode())
    print(f"Published release {tag}")


if __name__ == "__main__":
    main()
