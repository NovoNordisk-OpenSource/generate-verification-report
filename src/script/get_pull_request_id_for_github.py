import requests
import argparse
import json

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
        "Authorization": f"token {github_token}",
        "Accept": "application/vnd.github.v3+json"
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
                "closed_at": pr["closed_at"]
            })
        else:
            return "No pull request found for the given commit."
    else:
        return f"Failed to retrieve data. HTTP Status Code: {response.status_code}"

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
