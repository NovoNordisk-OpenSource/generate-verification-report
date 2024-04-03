import requests
import argparse
import json
from datetime import datetime

def format_timestamp(dt_string: str) -> str:
    # Convert to datetime object
    dt_object = datetime.strptime(dt_string, '%Y-%m-%dT%H:%M:%SZ')
    # Format datetime object as string
    formatted_string = dt_object.strftime('%Y-%m-%d %H:%M:%S')
    return formatted_string

def get_pull_request_details(commit_hash, github_token, repo):
    """
    Retrieve pull request details for a given commit hash from a GitHub repository.

    Parameters:
    - commit_hash: The hash of the commit to find the associated pull request.
    - github_token: Personal GitHub token for authentication.
    - repo: Repository name with owner (e.g., "octocat/Hello-World").

    Returns:
    A dictionary with pull request details or a message indicating no pull request found.
    """
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": f"Bearer {github_token}",
        "X-GitHub-Api-Version": "2022-11-28"
    }
    url = f"https://api.github.com/repos/{repo}/commits/{commit_hash}/pulls"

    response = requests.get(url, headers=headers)


    if response.status_code == 200:
        pull_requests = response.json()
        if pull_requests:
            pr = pull_requests[0]
            return json.dumps({
                "id": pr["id"],
                "number": pr["number"],
                "state": pr["state"],
                "title": pr["title"],
                "url": pr["html_url"],
                "closed_at": format_timestamp(pr["closed_at"])
            })
        else:
            return "No pull request found for the given commit."
    else:
        return f"Failed to retrieve data. HTTP Status Code: {response.status_code}. Error: {response.text}"

def main():
    parser = argparse.ArgumentParser(description="Get pull request details for a given commit hash from GitHub.")
    parser.add_argument("--commit", required=True, help="Commit hash")
    parser.add_argument("--token", required=True, help="GitHub Personal Access Token")
    parser.add_argument("--repo", required=True, help="Repository name with owner (e.g., octocat/Hello-World)")

    args = parser.parse_args()

    pr_details = get_pull_request_details(args.commit, args.token, args.repo)
    print(pr_details)

if __name__ == "__main__":
    main()
