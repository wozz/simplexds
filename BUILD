load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/wozz/simplexds
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = [
        "callbacks.go",
        "server.go",
        "store.go",
    ],
    importpath = "github.com/wozz/simplexds",
    visibility = ["//visibility:public"],
    deps = [
        "//mesh:go_default_library",
        "@com_github_envoyproxy_go_control_plane//envoy/api/v2:go_default_library",
        "@com_github_envoyproxy_go_control_plane//envoy/api/v2/core:go_default_library",
        "@com_github_envoyproxy_go_control_plane//envoy/service/discovery/v2:go_default_library",
        "@com_github_envoyproxy_go_control_plane//pkg/cache:go_default_library",
        "@com_github_envoyproxy_go_control_plane//pkg/log:go_default_library",
        "@com_github_envoyproxy_go_control_plane//pkg/server:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//peer:go_default_library",
    ],
)
