# Splitic

Splitic is an opinionated Go test runner for CI environments that can:

- Split tests for parallelism across different worker nodes based on timing information from previous runs.
- Output and merge junit test reports.
- Output and merge coverage reports.

## Usage

`splitic test` supports different timing information providers:

- **junit**

  `splitic test -provider junit --junit-filename junit.xml` allows you to specify a junit file to load
  from local disk for timing information.

- **gitlab**
  
  Using the [Test Report API](https://docs.gitlab.com/ee/api/pipelines.html#get-a-pipelines-test-report)
  GitLab will be used if `splitic` is run withing the GitLab CI environment. Alternatively, it can be configured 
  at the command line using `-provider gitlab`. Using `splitic test -provider --help` will show the various
  `-gitlab-*` options available to configure the provider (endpoint, branch to use for timing information etc.)

- **none**

  `splitic test -provider none` will not use timing information. Tests will be distributed evenly.

Tests are split deterministically into different buckets. The amount of buckets is equal to
`CI_NODE_TOTAL` or `-node-total`. The bucket chosen for any given worker node is specified through
`CI_NODE_INDEX` or `-node-index`.

