#!/bin/bash

if [ $# -eq 0 ]; then
  echo "Usage: link.sh url"
  exit 1
fi

# Escape ;
URL=`echo $1 | sed 's/;/%3b/g'`

printf '\033]1339;url='
echo -n $URL
printf '\a\n'

exit 0
