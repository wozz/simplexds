load("@io_bazel_rules_docker//go:image.bzl", "go_image")

go_image(
    name = "simplexds",
    embed = ["//cmd/simplexds:go_default_library"],
    visibility = ["//visibility:public"],
)
