# So far just contains tests for the C++ build rules.

# TODO(peterebden): now Python is becoming a plugin, we should choose something else to test with.
cc_shared_object(
    name = "so_test",
    srcs = ["so_test.cc"],
    out = "so_test.so",
    linker_flags = ["-bundle -undefined dynamic_lookup"] if is_platform(os = "darwin") else [],
    pkg_config_cflags = ["python3"],
    deps = [
        ":embedded_files",
    ],
)

# Make sure we have at least one working cc_binary rule.
cc_binary(
    name = "test_binary",
    srcs = ["test_binary.cc"],
    deps = [
        ":embedded_files",
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

# Tests on the C family of functions.
c_test(
    name = "c_test",
    srcs = ["test.c"],
)
