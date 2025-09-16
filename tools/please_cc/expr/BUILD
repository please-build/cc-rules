subinclude("///go//build_defs:go")

go_library(
    name = "expr",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = [
        "//tools/please_cc/...",
    ],
    deps = [
        "///third_party/go/github.com_PaesslerAG_gval//:gval",
        "///third_party/go/github.com_generikvault_gvalstrings//:gvalstrings",
        "//tools/please_cc/cctool",
    ],
)

go_test(
    name = "expr_test",
    srcs = glob(["*_test.go"]),
    deps = [
        "///third_party/go/github.com_stretchr_testify//assert",
        ":expr",
    ],
)
