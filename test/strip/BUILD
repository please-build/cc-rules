subinclude(
    "//build_defs:c",
    "//build_defs:cc",
)

package(
    cc = {
        # This is slightly contrived, but Mach-O executables always contain at least one symbol, even when stripped.
        # Emitting debug symbols ensures that Mach-O executables contain superfluous symbols that can later be stripped,
        # thus providing a clear distinction between stripped and unstripped executables for our tests.
        "default_opt_cflags": CONFIG.CC.DEFAULT_OPT_CFLAGS + ["-g"],
        "default_opt_cppflags": CONFIG.CC.DEFAULT_OPT_CPPFLAGS + ["-g"],
    },
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
        linker_flags = [f"-install_name @rpath/libc_so_{stripped}.so"] if CONFIG.OS == "darwin" else [],
        strip = strip,
    )
    data[f"c_{stripped}"] = c_binary(
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
        linker_flags = [f"-install_name @rpath/libcc_so_{stripped}.so"] if CONFIG.OS == "darwin" else [],
        strip = strip,
    )
    data[f"cc_{stripped}"] = cc_binary(
        name = f"cc_binary_{stripped}",
        srcs = ["binary.cpp"],
        linker_flags = binary_linker_flags + [f"-lcc_so_{stripped}"],
        strip = strip,
        deps = [f":cc_so_{stripped}"],
    )

# Ensure the executables still function as expected after being stripped.
test_cmd = ["for i in $DATA_C_STRIPPED $DATA_CC_STRIPPED $DATA_C_UNSTRIPPED $DATA_CC_UNSTRIPPED; do test $($i) = 42; done"]

if CONFIG.TARGET_OS == "darwin":
    test_cmd += [
        # `nm -a` should show that the unstripped executables contain debug symbols ("-")...
        "for i in $DATA_C_SO_UNSTRIPPED $DATA_C_UNSTRIPPED $DATA_CC_UNSTRIPPED; do nm -a $i | grep -qF ' - '; done",
        # ...but that the stripped executables don't.
        "for i in $DATA_C_SO_STRIPPED $DATA_C_STRIPPED $DATA_CC_STRIPPED; do ! nm -a $i | grep -qF ' - '; done",
    ]
else:
    test_cmd += [
        # `nm -a` should show that the unstripped executables contain at least one symbol...
        "for i in $DATA_C_SO_UNSTRIPPED $DATA_C_UNSTRIPPED $DATA_CC_UNSTRIPPED; do ! nm -a $i 2>&1 | grep -qF 'no symbols'; done",
        # ...but that the stripped executables contain none.
        "for i in $DATA_C_SO_STRIPPED $DATA_C_STRIPPED $DATA_CC_STRIPPED; do nm -a $i 2>&1 | grep -qF 'no symbols'; done",
    ]

gentest(
    name = "strip_test",
    test_cmd = test_cmd,
    data = data,
    exit_on_error = True,
    no_test_output = True,
)
