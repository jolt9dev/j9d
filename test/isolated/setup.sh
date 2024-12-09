#!/usr/bin/env bash

# get current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export SOPS_AGE_KEY_FILE="${DIR}/etc/keys.txt"