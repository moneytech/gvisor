load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

package(licenses = ["notice"])

# tests_bin must be buildable via `go build`, so it cannot have dependencies.
go_binary(
    name = "tests_bin",
    srcs = [
        "filter_input.go",
        "iptables_util.go",
        "main.go",
        "tests.go",
    ],
    visibility = ["//test/iptables:iptables_vis"],
)

# tests_bin must be buildable via `go build`, so it cannot have dependencies.
go_library(
    name = "tests",
    srcs = [
        "filter_input.go",  # keep
        "iptables_util.go",  # keep
        "tests.go",  # keep
    ],
    importpath = "gvisor.dev/gvisor/test/iptables/tests",
    visibility = ["//test/iptables:iptables_vis"],
)
