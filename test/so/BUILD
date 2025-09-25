subinclude("//build_defs:cc")

cc_shared_object(
    name = "libdolphin",
    srcs = ["dolphin.cpp"],
    hdrs = ["dolphin.hpp"],
)

cc_test(
    name = "so_test",
    srcs = ["so_test.cpp"],
    hdrs = ["dolphin.hpp"],
    linker_flags = [
        "-L" + package_name(),
        "-ldolphin",
        "-rpath '%s/%s'" % (
            ("@loader_path" if CONFIG.OS == "darwin" else "$ORIGIN"),
            package_name(),
        ),
    ],
    data = [":libdolphin"],
    deps = [":libdolphin"],
)
