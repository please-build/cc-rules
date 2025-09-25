subinclude("//build_defs:cc")

# cc-rules currently only supports the C++ module implementation in (mainline) Clang. Opt into these
# tests by enabling this plugin's "module_tests" feature flag.
if "module_tests" in CONFIG.CC.FEATURE_FLAGS:
    cc_module(
        name = "hello",
        srcs = ["hello.cc"],
        interfaces = ["hello.cppm"],
    )

    cc_binary(
        name = "modules_test",
        srcs = ["main.cc"],
        deps = [":hello"],
    )

    cc_module(
        name = "f1",
        srcs = ["f1.cc"],
        interfaces = ["f1.cppm"],
        deps = [":q1"],
    )

    cc_module(
        name = "q1",
        srcs = ["q1.cc"],
        interfaces = ["q1.cppm"],
        deps = [":f1"],
    )

    cc_test(
        name = "interdependent_modules_test",
        srcs = ["interdependent_modules_test.cc"],
        deps = [
            ":f1",
            ":q1",
        ],
    )
