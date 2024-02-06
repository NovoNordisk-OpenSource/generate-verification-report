[![Go Tests](https://github.com/BI-Data-Management-And-Analytics/verification-report-service/actions/workflows/on_commit.yml/badge.svg?branch=main)](https://github.com/BI-Data-Management-And-Analytics/verification-report-service/actions/workflows/on_commit.yml)

# Who is this for?
Software developers who understand the QMS toolchain.

# What does this do?
This service is responsible for generating verification reports given metadata about the QMS pipeline run, the test results, and Git context.

# How to use this?
> Mostly still TODO

## Installed prerequisites
- golang (version: >=1.21.6)
- dagger runtime (version: >= 0.9.7)

> The Dagger runtime is not required, but it renders nicer and logs are filtered properly when using it over:
> ```bash
> go run main.go
> ```

## Running the Dagger workflow
From within `/src`, run:

```bash
dagger run go run main.go
```

## Running tests
From within `/test`, run:

```bash
go test
```

# Generator workflow 
![](./doc/workflow.png)
