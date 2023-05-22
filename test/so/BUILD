# TODO(peterebden): now Python is becoming a plugin, we should choose something else to test with.
cc_shared_object(
    name = "shared",
    srcs = ["shared.cc"],
    out = "shared.so",
    linker_flags = ["-bundle -undefined dynamic_lookup"] if is_platform(os = "darwin") else [],
    deps = [
        "//test/embed:embedded_files",
    ],
)

cc_test(
    name = "shared_test",
    srcs = ["shared_test.cc"],
    data = [":shared"],
    deps = [
        "//test/embed:embedded_files",
    ],
)
