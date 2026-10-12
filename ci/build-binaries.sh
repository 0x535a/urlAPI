#!/bin/sh
set -eu

: "${CI_COMMIT_TAG:?This build requires a version tag}"
mkdir -p bin packages
for target_os in "$@"; do
  for target_arch in amd64 arm64; do
    ext=""
    ldflags=""
    case "${target_os}/${target_arch}" in
      linux/amd64) compiler=x86_64-linux-musl-gcc; ldflags=-static ;;
      linux/arm64) compiler=aarch64-linux-musl-gcc; ldflags=-static ;;
      windows/amd64) compiler=x86_64-w64-mingw32-clang; ext=.exe ;;
      windows/arm64) compiler=aarch64-w64-mingw32-clang; ext=.exe ;;
      darwin/amd64) compiler="clang -arch x86_64" ;;
      darwin/arm64) compiler="clang -arch arm64" ;;
      *) echo "Unsupported target: ${target_os}/${target_arch}" >&2; exit 1 ;;
    esac
    binary="urlAPI_${target_os}_${target_arch}${ext}"
    CGO_ENABLED=1 GOOS="$target_os" GOARCH="$target_arch" CC="$compiler" \
      CGO_LDFLAGS="$ldflags" go build -o "bin/$binary" .
    file "bin/$binary"
    package="urlAPI_${CI_COMMIT_TAG}_${target_os}_${target_arch}"
    mkdir -p "$package"
    cp "bin/$binary" README.md LICENSE "$package/"
    if [ "$target_os" = windows ]; then
      zip -r "packages/${package}.zip" "$package"
    else
      tar -czf "packages/${package}.tar.gz" "$package"
    fi
    rm -rf "$package"
  done
done
