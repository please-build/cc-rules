subinclude("///go//build_defs:go")

go_binary(
    name = "please_cc",
    srcs = ["please_cc.go"],
    deps = [
        "//tools/please_cc/cctool",
        "//tools/please_cc/expr",
    ],
    visibility = [
        "//package:all",
    ],
)

genrule(
    name = "version",
    srcs = ["VERSION"],
    outs = ["version.build_defs"],
    cmd = "echo PLEASE_CC_VERSION = \\\"$(cat $SRCS)\\\" > $OUTS",
    visibility = [
        "//package:all",
    ],
)
