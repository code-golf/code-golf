#!/bin/sh -e

export PATH=/usr/bin:/bin

[ "$1" = "--version" ] && exec /usr/local/bin/elm --version

cd /elm

# Inject holeArgs variable into code
# First delete the last line
sed -i '$d' src/Main.elm
shift
printf "holeArgs = [" >> src/Main.elm
while [[ $# -gt 0 ]]; do

  # Escape backslash, double quotes, and newlines
  escaped="$(\
    printf '%s' "$1" \
      | sed -e 's/\\/\\\\/g' \
      | sed -e 's/"/\\"/g' \
      | sed -e ':a' -e 'N' -e '$!ba' -e 's/\n/\\n/g' \
  )"

  printf '"%s"' "$escaped" >> src/Main.elm

  if [[ $# -gt 1 ]]; then
    printf ", " >> src/Main.elm
  fi

  shift
done
printf "]\n" >> src/Main.elm

# Write code to .elm file
cat - > src/M.elm

# Create output file
touch output

# Execute
set +e
err="$(/usr/local/bin/elm make src/Main.elm --optimize --output=elm.js 2>&1)"
status=$?
set -e

if [[ $status -eq 0 ]];then

  set +e
  err="$(node main.js > output)"
  status=$?
  set -e

  if [[ $status -eq 0 ]];then
    cat output
  else
    echo "$err" >&2
    exit $status
  fi

else
  echo "$err" >&2
  exit $status
fi
