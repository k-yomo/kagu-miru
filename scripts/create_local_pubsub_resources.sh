#!/bin/bash

set -euC

PROJECT=local
HOST=localhost:8085

function create_resources() {
  pubsub_cli create_subscription item-update item-update.elasticsearch-indexer --create-if-not-exist -p $PROJECT -h $HOST
  pubsub_cli create_subscription item-update item-update.item-spanner-inserter --create-if-not-exist -p $PROJECT -h $HOST
}

NEXT_WAIT_TIME=0
until create_resources || [ $NEXT_WAIT_TIME -eq 10 ]; do
   sleep $(( NEXT_WAIT_TIME++ ))
done
