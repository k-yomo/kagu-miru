#!/bin/bash

set -euC

function create_index_if_not_exist() {
  local index_name="$1"
  local status=$(curl -s "localhost:9200/$index_name/" | jq '.status')
  if [ "$status" = '404' ]; then
    curl -XPUT  -H "Content-Type: application/json" -d @./defs/elasticsearch/mappings/$index_name.json "localhost:9200/$index_name/"
    echo
  fi
}


function main() {
  create_index_if_not_exist 'items'
  create_index_if_not_exist 'items.query_suggestions'
}

main
