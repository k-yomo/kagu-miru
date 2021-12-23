#!/bin/bash

set -euC

function create_index_if_not_exist() {
  local filepath="$1"
  filename=$(basename $filepath)
  index_name=${filename%.json}
  local status=$(curl -s "localhost:9200/$index_name/" | jq '.status')
  if [ "$status" = '404' ]; then
    curl -XPUT  -H "Content-Type: application/json" -d @$filepath "localhost:9200/$index_name/"
    echo # for line break
  fi
}


function main() {
  for filepath in defs/elasticsearch/mappings/*.json; do
    create_index_if_not_exist $filepath
  done
}

main
