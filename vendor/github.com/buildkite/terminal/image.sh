#!/bin/bash

# externalimg.sh url
if [ $# -eq 0 ]; then
  echo "Usage: image.sh url"
  exit 1
fi

printf '\033]1338;url='
echo -n $1
printf '\a\n'

exit 0
