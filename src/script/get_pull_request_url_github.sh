#!/bin/bash

if [[ -n $pr ]]; then
    # Code to execute if $pr is set and has a value
    prUrl=$(jq -r .url <<< $pr)
    echo $prUrl
else
    # Code to execute if $pr is not set or has no value
    echo "'pr' variable is not set or has no value"
    exit 1
fi