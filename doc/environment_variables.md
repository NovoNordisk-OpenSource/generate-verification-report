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

In some cases it may be necessary to override the default values of these environment variables. This can be done by [setting override environment variables](#overriding-default-github-environment-variables) that the Dagger workflow will look for.

### Setting default environment variables

> This command assumes running on linux or macOS.

You can set the default environment variables manually, e.g. for local use, by running the following command:

```shell
export GITHUB_REPOSITORY=nn-dma/generate-verification-report-test
export GITHUB_REF_NAME=main
export GITHUB_SHA=724a0a893e760ae2df3f809985ee55feda4cb7a9
export GITHUB_RUN_ID=0123456789
export GITHUB_TOKEN=ghp_1234567890abcdefghijklmnopqrstuvwxyz
```

See the section [default environment variables](#default-environment-variables) for more information on what each variable is used for.

### Default environment variables

> These commands assume running on linux or macOS.

#### GITHUB_REPOSITORY

These are related to the runtime context and will be set automatically by the GitHub workflow. If you run the Dagger workflow locally or have a reason to override contextual values, they can be manually set and overridden.

It is used to generate links and looking up information for generating the verification report.

Example:

```shell
export GITHUB_REPOSITORY=nn-dma/generate-verification-report-test
```

#### GITHUB_REF_NAME

The `GITHUB_REF_NAME` environment variable is expected to be the name of the branch or tag that triggered the workflow. 

It is used when looking up information while generating the verification report.

Example:

```shell
export GITHUB_REF_NAME=main
```

#### GITHUB_SHA

The `GITHUB_SHA` environment variable is expected to be the commit SHA that triggered the workflow. 

It is used to generate links to commits in the verification report.

Example:

```shell
export GITHUB_SHA=724a0a893e760ae2df3f809985ee55feda4cb7a9
```

#### GITHUB_RUN_ID

The `GITHUB_RUN_ID` environment variable is expected to be the unique identifier of the workflow run. 

It is used to generate links to the workflow run in the verification report.

Example:

```shell
export GITHUB_RUN_ID=123456789
```

#### GITHUB_TOKEN

This grants access to the GitHub API and is used to fetch information about the repository, issues, pull request and the workflow run. It is possible to run the workflow against a repository in another place, which is the most common reason—aside from running locally—to override this value.

Example:

```shell
export GITHUB_TOKEN=ghp_1234567890abcdefghijklmnopqrstuvwxyz
```

## Overriding default GitHub environment variables

You may want to override the default values of the GitHub environment variables. This can be done by setting the following override environment variables that the Dagger workflow will look for.

### Setting override environment variables

> This command assumes running on linux or macOS.

You can set the override environment variables manually, e.g. for local use, by running the following command:

```shell
export OVERRIDE_GITHUB_REPOSITORY=nn-dma/generate-verification-report-test
export OVERRIDE_GITHUB_REF_NAME=main
export OVERRIDE_GITHUB_SHA=724a0a893e760ae2df3f809985ee55feda4cb7a9
```

See the section [override environment variables](#override-environment-variables) for more information on what each variable is used for.

### Override environment variables

> These commands assume running on linux or macOS.

#### OVERRIDE_GITHUB_REPOSITORY

This environment variable is used to override the default value of `GITHUB_REPOSITORY`.

Example:

```shell
export OVERRIDE_GITHUB_REPOSITORY=nn-dma/generate-verification-report-test
```

#### OVERRIDE_GITHUB_REF_NAME

This environment variable is used to override the default value of `GITHUB_REF_NAME`.

Example:

```shell
export OVERRIDE_GITHUB_REF_NAME=main
```

#### OVERRIDE_GITHUB_SHA

This environment variable is used to override the default value of `GITHUB_SHA`. 

It is useful when you want to generate a verification report for a specific commit or pull request that did not necessarily trigger the latest workflow run, or if you want to generate the verification report based on a pull request in another repository.

Example:

```shell
export OVERRIDE_GITHUB_SHA=724a0a893e760ae2df3f809985ee55feda4cb7a9
```

## Unsetting all environment variables

> This command assumes running on linux or macOS.

You can unset the all possible environment variables manually by running:

```shell
unset OVERRIDE_GITHUB_REPOSITORY
unset OVERRIDE_GITHUB_REF_NAME
unset OVERRIDE_GITHUB_SHA
unset GITHUB_REPOSITORY
unset GITHUB_REF_NAME
unset GITHUB_SHA
unset GITHUB_RUN_ID
unset GITHUB_TOKEN
```