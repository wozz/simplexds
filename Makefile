.PHONY: run update update-repos docker

run: 
	bazelisk run //cmd/simplexds:simplexds

update:
	bazelisk run :gazelle

update-repos:
	bazelisk run :gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories
	bazelisk run :gazelle

docker:
	bazelisk build //docker:simplexds --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
