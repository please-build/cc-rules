subinclude("//build_defs:cc")

cc_binary(
    name = "test_binary",
    srcs = ["test_binary.cc"],
)

gentest(
    name = "cc_binary_test",
    data = [":test_binary"],
    labels = [
        "cc",
    ],
    no_test_output = True,
    test_cmd = "$(exe :test_binary)",
)
