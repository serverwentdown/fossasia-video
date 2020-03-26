#!/bin/sh

set -e

mode=loudnorm

file="$1"
outfile="${file%.*}.audiotweaked.${file##*.}"

if [[ -z "$file" ]]; then
	echo "Usage: $0 videofile"
	exit 1
fi

if [[ "$mode" == "mean" ]]; then
	# A simple way to detect mean_volume and use it to tweak the entire file
	ffmpeg -i "$file" -filter:a volumedetect -f null /dev/null
	read -p "Enter mean_volume: " volume
	ffmpeg -i "$file" -filter:a "volume=$volume" -af "highpass=f=200, lowpass=f=3000" "$outfile"
elif [[ "$mode" == "loudnorm" ]]; then
	# Use the recommended loudnorm filter
	ffmpeg -i "$file" -vcodec copy -af "loudnorm" "$outfile"
else
	echo "Unknown mode"
	exit 1
fi
