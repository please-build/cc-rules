cc_binary(
    name = "test_binary",
    srcs = ["test_binary.cc"],
    deps = [
        "//test/embed:embedded_files",
    ],
)

gentest(
    name = "cc_binary_test",
    data = [":test_binary"],
    labels = [
        "cc",
        "embed",
    ],
    no_test_output = True,
    test_cmd = "$(exe :test_binary)",
)
