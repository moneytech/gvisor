# iptables Tests

iptables rules are tested via the Dockerfile in this directory.

## Test Structure

Each test implements `Testcase`, and provides (1) a function to run inside the
container and (2) a function to run locally. Those processes are given each
others' IP addresses. The test succeeds when both functions succeed.

The function inside the container (`ContainerAction`) typically sets some
iptables rules and then tries to send or receive packets. The local function
(`LocalAction`) will typically just send or receive packets.

### Adding Tests

1) Add your test to the `tests` directory.

2) Add it to the list of tests in `tests/tests.go`.

3) Add it to `iptables_test.go` (see the other tests in that file).

Your test is now runnable with bazel (see below for instructions on how to run
all tests or only the new one).

## Setup

1) [Install and configure Docker](https://docs.docker.com/install/)

2) [Install and configure Bazel](https://bazel.build/)

3) Build the Docker container:

```bash
$ docker build -t iptables-tests .
```

## Running the Tests

This will run the iptables test both in a standard and runsc docker container:

```bash
$ bazel test //test/iptables_test_{runc,runsc}
```

To run a single test:

```bash
$ bazel test //test/iptables_test_{runc,runsc} --test_filter=<TESTNAME>
```
