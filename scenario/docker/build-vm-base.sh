#!/bin/bash
set -euo pipefail

architecture="$1"
checkout="${AZL_IMAGE_BUILD_CHECKOUT:-/azurelinux}"
artifact_directory="${AZL_IMAGE_BUILD_ARTIFACTS:-/artifacts}"
output_directory="$checkout/base/out/images/vm-base"
artifact="$artifact_directory/vm-base.qcow2"

if [[ -e "$artifact" || -L "$artifact" ]]; then
    echo "vm-base.qcow2 already exists" >&2
    exit 1
fi

cd "$checkout"
azldev image build vm-base --arch "$architecture"

shopt -s nullglob
images=("$output_directory"/*.raw)
if [[ "${#images[@]}" -ne 1 || ! -f "${images[0]}" ]]; then
    echo "expected exactly one vm-base raw image" >&2
    exit 1
fi

temporary="$(mktemp "$artifact_directory/.vm-base.qcow2.XXXXXX")"
trap 'rm -f "$temporary"' EXIT

qemu-img convert -f raw -O qcow2 -c "${images[0]}" "$temporary"
qemu-img check "$temporary"
ln -T -- "$temporary" "$artifact"

rm -f "$temporary"
trap - EXIT
