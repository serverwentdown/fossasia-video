#!/bin/sh

set -e

temp=/tmp/lossless-cut
download=https://github.com/mifi/lossless-cut/releases/download/v3.17.4/LosslessCut-linux.tar.bz2

mkdir -p "$temp"

if [[ ! -f "$temp/lossless-cut" ]]; then
	echo "lossless-cut not found in $temp"
	echo "Downloading lossless-cut..."
	wget -nv --show-progress "$download" -O "$temp/archive.tar.bz2"
	echo "Extracting lossless-cut..."
	tar -xf "$temp/archive.tar.bz2" -C "$temp" --strip-components=1
	echo "Successfully obtained lossless-cut"
fi

"$temp/lossless-cut" "$@"



