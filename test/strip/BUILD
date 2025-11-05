subinclude(
    "//build_defs:c",
    "//build_defs:cc",
)

binary_linker_flags = [
    "-L" + package_name(),
    "-rpath '%s'" % ("@executable_path" if CONFIG.OS == "darwin" else "$ORIGIN"),
]

data = {}

for strip in [True, False]:
    stripped = "stripped" if strip else "unstripped"

    data[f"c_so_{stripped}"] = c_shared_object(
        name = f"c_so_{stripped}",
        srcs = ["so.c"],
        out = f"libc_so_{stripped}.so",
        hdrs = ["so.h"],
        linker_flags = [f"-install_name @rpath/lib_c_{stripped}.so"] if CONFIG.OS == "darwin" else [],
        strip = strip,
    )
    data[f"c_binary_{stripped}"] = c_binary(
        name = f"c_binary_{stripped}",
        srcs = ["binary.c"],
        hdrs = ["so.h"],
        linker_flags = binary_linker_flags + [f"-lc_so_{stripped}"],
        strip = strip,
        deps = [f":c_so_{stripped}"],
    )
    data[f"cc_so_{stripped}"] = cc_shared_object(
        name = f"cc_so_{stripped}",
        srcs = ["so.cpp"],
        out = f"libcc_so_{stripped}.so",
        hdrs = ["so.hpp"],
        linker_flags = [f"-install_name @rpath/lib_cc_{stripped}.so"] if CONFIG.OS == "darwin" else [],
        strip = strip,
    )
    data[f"cc_binary_{stripped}"] = cc_binary(
        name = f"cc_binary_{stripped}",
        srcs = ["binary.cpp"],
        linker_flags = binary_linker_flags + [f"-lcc_so_{stripped}"],
        strip = strip,
        deps = [f":cc_so_{stripped}"],
    )

gentest(
    name = "strip_test",
    test_cmd = [
        "for i in $DATA_C_BINARY_STRIPPED $DATA_CC_BINARY_STRIPPED $DATA_C_BINARY_UNSTRIPPED $DATA_CC_BINARY_UNSTRIPPED; do test $($i) = 42; done",
        "for i in $DATA_C_SO_UNSTRIPPED $DATA_C_BINARY_UNSTRIPPED $DATA_CC_BINARY_UNSTRIPPED; do file $i | grep -q 'not stripped'; done",
        "for i in $DATA_C_SO_STRIPPED $DATA_C_BINARY_STRIPPED $DATA_CC_BINARY_STRIPPED; do file $i | grep -q ', stripped'; done",
    ],
    data = data,
    exit_on_error = True,
    no_test_output = True,
)
