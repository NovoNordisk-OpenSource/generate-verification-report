[![Unit Tests](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit.yml/badge.svg?branch=main)](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit.yml) [![Integration Tests](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit_workflow.yml/badge.svg?branch=main)](https://github.com/nn-dma/generate-verification-report/actions/workflows/on_commit_workflow.yml)

# Who is this for?
Software developers who understand the QMS toolchain.

# What does this do?
This workflow is responsible for generating verification reports given metadata about the QMS pipeline run, the test results, and Git context.
It is intended to be called from within the GitHub action that generates a verification report as part of the QMS pipeline.

However, it can technically be invoked from anywhere capable of running Dagger and Docker Engine.

# How to use this?
First, go to the `/src` directory.

Parameters are provided here by editing the `parameters.json` file.

Logs will be written to a `run.log` file in the same place everytime the Dagger workflow is run. This file is appended between runs.

Output from any run will be placed in a directory here called `output`. The content in this will be the verification report HTML file.
Any output in this directory will be overwritten between runs unless the generated filename is different because of different intput parameters.

> TODO: Determine how test results are ingested—probably from an `input` directory.

## Installed prerequisites
- golang (version: >=1.22.0)
- dagger runtime (version: >= 0.9.10)

> The Dagger runtime is not required, but it renders nicer and logs are filtered properly when using it over:
> ```text
> go run main.go
> ```

An additional requirement is that the executing host environment can reach the public Docker Hub image registry.

## Running the Dagger workflow
The production codebase lives in the `/src` directory.

From within `/src`, run:

```text
dagger run go run main.go
```

## Running tests
The test codebase lives in the `/test` directory and is consuming the production codebase. It is thus not part of the production codebase.

Two subgroups of tests exist: `unit tests` and `integration tests`. They have their own folders in `/test` and are run separately.

From within `/test/integration` or `/test/unit`, run:

```text
go test
```

# Generator workflow 
![](./doc/workflow.png)
