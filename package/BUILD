subinclude(
    "//build_defs:arch",
    "//tools/please_cc:version",
)

tools = {
    "please_cc": PLEASE_CC_VERSION,
}

target_names = {}

def cross_compile(tool:str, version:str, arch:str):
    target = target_names.get(tool, tool)
    return build_rule(
        name = f"{tool}_{arch}",
        srcs = [f"///{arch}//tools/{tool}:{target}"],
        outs = [f"{tool}-{version}-{arch}"],
        cmd = "mv $SRCS $OUTS",
        labels = ["manual"],
    )

for tool, version in tools.items():
    filegroup(
        name = f"{tool}_release_files",
        srcs = [cross_compile(tool, version, arch) for arch in SUPPORTED_ARCHITECTURES],
        labels = [
            f"hlink:plz-out/package/{tool}",
            "manual",
        ],
    )
