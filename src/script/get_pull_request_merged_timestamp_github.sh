#!/bin/bash

if [[ -n $pr ]]; then
    # Code to execute if $pr is set and has a value
    prMergedTimestamp=$(jq -r .merged_at <<< $pr)
    echo $prMergedTimestamp
else
    # Code to execute if $pr is not set or has no value
    echo "'pr' variable is not set or has no value"
    #exit 1
fi