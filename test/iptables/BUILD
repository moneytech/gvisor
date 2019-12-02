load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(licenses = ["notice"])

go_library(
    name = "iptables",
    srcs = ["iptables.go"],
    importpath = "gvisor.dev/gvisor/test/iptables",
)

go_test(
    name = "iptables_test_runc",
    srcs = ["iptables_test.go"],
    args = ["-runtime=runc"],
    embed = [":iptables"],
    tags = [
        # Requires docker and runsc to be configured before the test runs.
        "manual",
        # Don't include test target in wildcard target patterns.
        "local",
    ],
    deps = [
        "//pkg/log",
        "//runsc/dockerutil",
        "//test/iptables/tests",
    ],
)

go_test(
    name = "iptables_test_runsc",
    srcs = ["iptables_test.go"],
    args = ["-runtime=runsc"],
    embed = [":iptables"],
    tags = [
        # Requires docker and runsc to be configured before the test runs.
        "manual",
        # Don't include test target in wildcard target patterns.
        "local",
    ],
    deps = [
        "//pkg/log",
        "//runsc/dockerutil",
        "//test/iptables/tests",
    ],
)

package_group(
    name = "iptables_vis",
    packages = [
        "//test/iptables/...",
    ],
)
