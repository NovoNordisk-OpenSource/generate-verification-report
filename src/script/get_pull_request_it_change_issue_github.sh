#!/bin/bash

if [[ -n $pr ]]; then
    # Code to execute if $pr is set and has a value
    prItChangeIssueNumber=$(jq -r '. | select(.title != null) | .title | match("(IT Change) #\\d+"; "g") | .string | sub(".+ #";"")' <<< $pr)
    echo $prItChangeIssueNumber
else
    # Code to execute if $pr is not set or has no value
    echo "'pr' variable is not set or has no value"
    #exit 1
fi