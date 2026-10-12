# Source this script in the Linux build job so PATH persists for the build.
set -eu

cross_root="${CI_PROJECT_DIR:-$PWD}/.cross"
mkdir -p "$cross_root"

curl --http1.1 --fail --silent --show-error --location --retry 3 \
  https://musl.cc/aarch64-linux-musl-cross.tgz \
  --output "$cross_root/aarch64-linux-musl-cross.tgz"
printf '%s  %s\n' c909817856d6ceda86aa510894fa3527eac7989f0ef6e87b5721c58737a06c38 \
  "$cross_root/aarch64-linux-musl-cross.tgz" | sha256sum -c -
tar -xzf "$cross_root/aarch64-linux-musl-cross.tgz" -C "$cross_root"

llvm_mingw=llvm-mingw-20250430-ucrt-ubuntu-22.04-x86_64
curl --http1.1 --fail --silent --show-error --location --retry 3 \
  "https://github.com/mstorsjo/llvm-mingw/releases/download/20250430/${llvm_mingw}.tar.xz" \
  --output "$cross_root/llvm-mingw.tar.xz"
tar -xf "$cross_root/llvm-mingw.tar.xz" -C "$cross_root"
rm "$cross_root/aarch64-linux-musl-cross.tgz" "$cross_root/llvm-mingw.tar.xz"

export PATH="$cross_root/aarch64-linux-musl-cross/bin:$cross_root/$llvm_mingw/bin:$PATH"
