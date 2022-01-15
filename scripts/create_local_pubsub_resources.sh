#!/bin/bash

set -euC

PROJECT=local
HOST=localhost:8085

function create_resources() {
  pubsub_cli create_subscription item-update item-update.item-indexer --create-if-not-exist -p $PROJECT -h $HOST
}

NEXT_WAIT_TIME=0
until create_resources || [ $NEXT_WAIT_TIME -eq 10 ]; do
   sleep $(( NEXT_WAIT_TIME++ ))
done
