#!/bin/bash

# This script will evaluate all `plot.star` files with warplark.
#
# NOTE: This is not the intended way to use warplark. You likely want to use
# `warpforge plan generate ./...` instead. This script is only for testing
# warplark without using Warpforge.

find . -name plot.star | while IFS= read -r in_file
do 
    out_file="$(dirname $in_file)/plot.wf"
    echo "$in_file -> $out_file"
    warplark $in_file > $out_file
done
