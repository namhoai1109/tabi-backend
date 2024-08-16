#!/bin/bash

cd functions
ARGS="$@"

echo "$ sls ${ARGS}"
echo ""

sls "$@"