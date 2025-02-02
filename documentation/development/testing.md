# Testing

Tests are written in [Ruby](https://www.ruby-lang.org) for historical reasons
and because that allows running them
[in parallel](https://github.com/grosser/parallel_tests).

The end-to-end tests are located in [features](../../features) and written in
[Cucumber](https://github.com/cucumber/cucumber-ruby). Unit tests are written as
normal Go tests using [Ginkgo](https://github.com/onsi/ginkgo).

## Running Tests

```bash
# running the different test types
make test       # runs all tests
make lint       # runs the linters
make lint-go    # runs the Go linters
make cuke       # runs the feature tests

# running individual scenarios/features
make cuke dir=<path>

# running tests in parallel
make cuke [cucumber parameters]
# set the environment variable PARALLEL_TEST_PROCESSORS to override the
# auto-detected number of processors

# auto-fixing formatting issues
make fix
```

Git Town's [CI server](https://circleci.com/gh/Originate/git-town) automatically
tests all commits and pull requests, and notifies you via email and through
status badges in pull requests about problems.

## Debugging

**See the CLI output of Ruby specs:** set the `DEBUG_COMMANDS` environment
variable while running your specs:

```bash
$ DEBUG_COMMANDS=true cucumber <filename>[:<lineno>]
```

Alternatively, you can also add a `@debug-commands` flag to the respective
Cucumber spec:

```cucumber
@debug-commands
Scenario: foo bar baz
  Given ...
```

For even more detailed output, you can use the `DEBUG` variable or tag in a
similar fashion. If set, Git Town prints every shell command executed during the
tests (includes setup, inspection of the Git status, and the Git commands), and
the respective console output.

**See the CLI output of Go specs:**

- add a tag `@debug` to see the output of all shell commands

**debug a Godog Cucumber spec in VSCode:**

- open `main_test.go`
- in this file, change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

## Mocking

Certain tests require the Git remote to be set to a real value on GitHub or
Bitbucket. This causes `git push` operations to go to GitHub during testing,
which is undesirable. To circumvent this problem, Git Town allows to mock the
Git remote by setting the Git configuration value `git-town.testing.remote-url`
to the respective value. To keep this behavior clean and secure, this also
requires an environment variable `GIT_TOWN_ENV` to be set to `test`.

## Auto-running tests

The Git Town code base works with
[Tertestrial](https://github.com/Originate/tertestrial-server).
