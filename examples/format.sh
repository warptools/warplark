#!/bin/bash

# this script formats the example lark files
#
# it requires yapf: `pip install yapf`
# it must be run from the examples dir
#
# yapf won't search for non-python files, so we must collect the .star
# files using find

find $(dirname $0) -type f -name "*.star" | xargs yapf -i
