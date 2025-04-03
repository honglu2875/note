load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "note_lib",
    srcs = ["main.go"],
    importpath = "github.com/honglu2875/note/note",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "note",
    embed = [":note_lib"],
    visibility = ["//visibility:public"],
)
