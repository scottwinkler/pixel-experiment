#!/bin/bash
for f in $(find . -name '*.wav'); do
 name=$(basename $f | cut -f 1 -d '.')
 dir=$(dirname $f)
 ffmpeg -y -i $f -c:a flac -sample_fmt s16  $dir/$name.flac;
done