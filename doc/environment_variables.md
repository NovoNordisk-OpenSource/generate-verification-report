# Environment variables

> This document describes functionality related to executing the Dagger workflow in a GitHub context only.

## How they are used

The Dagger workflow expects the following environment variables to be set.

If any of these are missing or have invalid values, the workflow will fail. Some validation is in place within the Dagger workflow, but it is limited to checking for the presence of the variables.

```shell
GITHUB_REPOSITORY
GITHUB_REF_NAME
GITHUB_SHA
GITHUB_RUN_ID
GITHUB_TOKEN
```

When the Dagger workflow is run in a GitHub context, these environment variables are set automatically by the GitHub workflow. If you run the Dagger workflow locally or have a reason to override contextual values, they can be manually set and overridden.

The variables are expected to follow the official GitHub [documentation](https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables).

The `GITHUB_TOKEN` is expected to be the dynamically generated GitHub token during workflow run, but it can also be a personal access token (PAT) with the necessary permissions.

### GITHUB_REPOSITORY

These are related to the runtime context and will be set automatically by the GitHub workflow. If you run the Dagger workflow locally or have a reason to override contextual values, they can be manually set and overridden.

Example:

```shell
GITHUB_REPOSITORY=nn-dma/generate-verification-report-test
```

### GITHUB_REF_NAME

The `GITHUB_REF_NAME` environment variable is expected to be the name of the branch or tag that triggered the workflow.

Example:

```shell
GITHUB_REF_NAME=main
```

### GITHUB_SHA

The `GITHUB_SHA` environment variable is expected to be the commit SHA that triggered the workflow.

Example:

```shell
GITHUB_SHA=724a0a893e760ae2df3f809985ee55feda4cb7a9
```

### GITHUB_RUN_ID

The `GITHUB_RUN_ID` environment variable is expected to be the unique identifier of the workflow run.

Example:

```shell
GITHUB_RUN_ID=123456789
```

### GITHUB_TOKEN

Example:

```shell
GITHUB_TOKEN=ghp_1234567890abcdefghijklmnopqrstuvwxyz
```