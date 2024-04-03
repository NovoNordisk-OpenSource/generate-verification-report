#!/bin/bash

commit_hash=$1
github_token=$2
repo=$3

self_dir=$(dirname "$(realpath "$0")")
#echo "self directory: $self_dir"

python_file="get_pull_request_details_for_hash_github.py"
python_script_path="$self_dir/$python_file"
#echo "python script path: $python_script_path"

pr=$(python3 "$python_script_path" --commit $commit_hash --token $github_token --repo $repo)
prUrl=$(jq -r .url <<< $pr)
echo $prUrl