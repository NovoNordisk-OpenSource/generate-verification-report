[![Unit Tests](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit.yml/badge.svg?branch=main)](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit.yml) [![Integration Tests](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit_workflow.yml/badge.svg?branch=main)](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit_workflow.yml)

# Who is this for?
Software developers who understand the QMS toolchain.

# What does this do?
This workflow is responsible for generating verification reports given metadata about the QMS pipeline run, the test results, and Git context.
It is intended to be called from within the GitHub action that generates a verification report as part of the QMS pipeline.

However, it can technically be invoked from anywhere capable of running Dagger and Docker Engine.

![](./doc/dagger_workflow.png)

# How to use this?
First, go to the `/src` directory. Following paths mentioned are relative to this.

### Inputs

Parameters are provided here by editing the `parameters.json` file. It is located in the `/input` directory *(not to be confused with the **`inputs`** directory, which is a Go package)*.

Test results provided as input to the Dagger worklow must be placed in the `/input/testresults` directory. For now, the test results must be in the form of one JSON file per test case result and each must be in the Allure-format.

### Outputs

Output from any run will be placed in a directory here called `/output`. The content in this will be the verification report HTML file.
Any output in this directory will be overwritten between runs unless the generated filename is different because of different intput parameters.

#### Logs

Logs will be written to a `run.log` file in the same place everytime the Dagger workflow is run. This file is appended between runs.

## Installed prerequisites
- golang (version: >=1.22.1)
- dagger runtime (version: >= 0.10.1)

An additional requirement is that the executing host environment can reach the public Docker Hub image registry.

## Running the Dagger workflow
The production codebase lives in the `/` directory.

From within `/`, run:

```text
dagger run go run -C src main.go
```

## Running tests
The test codebase lives in the `/test` directory and is consuming the production codebase. It is thus not part of the production codebase.

Two subgroups of tests exist: `unit tests` and `integration tests`. They have their own folders in `/test` and are run separately.

To execute unit tests , from within `/` run:

```text
go test -C test/unit -v
```

To execute integration tests , from within `/` run:

```text
go test -C test/integration -v
```

## Running workflows locally with `act`
The GitHub workflows can be executed locally with [act](https://github.com/nektos/act). Install with Homebrew or another package manager. Using `act` also requires Docker Engine to be installed.

When running the workflows locally, `act` might initially ask you which container size to use. *Medium* should work fine for the time being.

To execute the **unit tests** workflow, from within `/`, run:
```text
act --container-architecture linux/amd64 -W .github/workflows/on_commit.yml
```

To execute the **integration tests** workflow, from within `/`, run:
```text
act --container-architecture linux/amd64 -W .github/workflows/on_commit_workflow.yml
```
