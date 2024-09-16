# Uploading artifacts

You will need to provide two artifacts for the verification report generation workflow to run successfully:

1. The repository under test, named `repository` by default.
2. The test results for two environments:
   a. Validation environment, named `testresults-validation` by default.
   b. Production environment, named `testresults-production` by default.

While it is possible to change the names of the artifacts, there is really no reason to do so as it only creates work for you to maintain these new names. The above are the _sensible defaults_.

## How to upload using GitHub Actions

It is recommended to use the `actions/upload-artifact` action to upload the artifacts. However, be aware that you need to specify that hidden files must be included in order to upload the `.git` directory for the `repository` artifact.

Like this:

```yaml
- name: Upload repository
  uses: actions/upload-artifact@v4.4.0
  with:
    name: repository
    path: input/repository
    include-hidden-files: true
```
