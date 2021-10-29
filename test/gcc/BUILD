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

# Tests on the C family of functions.
c_test(
    name = "c_test",
    srcs = ["test.c"],
)
