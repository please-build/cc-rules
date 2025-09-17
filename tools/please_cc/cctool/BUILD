subinclude("///go//build_defs:go")

go_library(
    name = "cctool",
    srcs = glob(
        ["*.go"],
        exclude = [
            "*_test.go",
            "generate_test_data.go",
        ],
    ),
    visibility = [
        "//tools/please_cc/...",
    ],
)

go_test(
    name = "cctool_test",
    srcs = glob(["*_test.go"]),
    data = ["test_data"],
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        ":cctool",
    ],
)

go_binary(
    name = "generate_test_data",
    srcs = ["generate_test_data.go"],
    deps = [
        "///third_party/go/github.com_thought-machine_go-flags//:go-flags",
    ],
)
