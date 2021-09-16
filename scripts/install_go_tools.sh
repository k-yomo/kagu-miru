cat ./tools/tools.go | awk -F'"' '/_/ {print $2}' | xargs -tI {} go install {}
