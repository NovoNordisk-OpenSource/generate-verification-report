#!/bin/bash

if [[ -n $pr ]]; then
    # Code to execute if $pr is set and has a value
    prItChangeIssueNumber=$(jq -r '. | select(.title != null) | .title | match("(IT Change) #\\d+"; "g") | .string | sub(".+ #";"")' <<< $pr)
    # Check if the command was successful
    if [ $? -ne 0 ]; then
        echo "error: could not extract 'IT Change' issue number from pull request details"
        echo "given: '$pr'"
        #exit 1
    else
        echo $prItChangeIssueNumber
    fi
else
    # Code to execute if $pr is not set or has no value
    echo "error: 'pr' variable is not set or has no value"
    #exit 1
fi