#!/bin/bash

cd functions
ARGS="$@"

sls deploy

# if ARG = involk then run sls invoke
if [ "$ARGS" == "invoke" ]; then
    sls invoke --function Migration
    exit 0
fi