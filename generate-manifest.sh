#!/bin/bash

set -euo pipefail

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <BUILD_DIR> <VERSION> <PROVIDER_NAME>" >&2
    exit 1
fi

BUILD_DIR="$1"
VERSION="$2"
PROVIDER="$3"

if [ ! -d "$BUILD_DIR" ]; then
  echo "Error: Build directory '$BUILD_DIR' not found." >&2
  exit 1
fi

cd "$BUILD_DIR"

VERSION_VALUE="${VERSION#v}"
MANIFEST="terraform-provider-${PROVIDER}_${VERSION_VALUE}_manifest.json"
PLATFORMS_JSON="["
first=1

# Use nullglob to prevent errors if no files match
shopt -s nullglob
files=(terraform-provider-${PROVIDER}_v${VERSION_VALUE}_*.zip)
shopt -u nullglob


if [ ${#files[@]} -eq 0 ]; then
    echo "No provider zip files found for version ${VERSION_VALUE}. Generating empty manifest." >&2
    jq -n --arg version "$VERSION_VALUE" --argjson platforms "[]" '{version:$version, protocols:["5.0"], platforms:$platforms}' > "$MANIFEST"
    exit 0
fi

for f in "${files[@]}"; do
  sha=$(shasum -a 256 "$f" | awk '{print $1}')
  base=${f%.zip}
  os_arch=${base#terraform-provider-${PROVIDER}_${VERSION_VALUE}_}
  os=${os_arch%_*}
  arch=${os_arch#${os}_}
  entry=$(jq -n --arg os "$os" --arg arch "$arch" --arg filename "$f" --arg shasum "$sha" '{os:$os, arch:$arch, filename:$filename, shasum:$shasum}')
  if [ $first -eq 1 ]; then PLATFORMS_JSON="[$entry"; first=0; else PLATFORMS_JSON="$PLATFORMS_JSON,$entry"; fi
done

PLATFORMS_JSON="$PLATFORMS_JSON]"
jq -n --arg version "$VERSION_VALUE" --argjson platforms "$PLATFORMS_JSON" '{version:$version, protocols:["5.0"], platforms:$platforms}' > "$MANIFEST"

echo "Manifest generated: $MANIFEST"