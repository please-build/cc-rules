tools = {
    "please_cc": {
        "version": "0.1.1",
        "hashes": {
            "darwin_amd64": "b0f62ee03e2d685ef6b6e3859c61dc452ad4bd361986676e8b8bdac94ddd0ed3",
            "darwin_arm64": "f62f9824ec5b3ec3b8f6d07968d1851994731d752a295994868c66e56bbe0966",
            "freebsd_amd64": "47a6087bde791b66a77eab3ca27fa7db4215c0e8a0419822f38e70dc7c5740bb",
            "linux_amd64": "a44fb9f24097cca844f69f344357fdeeacd94e40d3cf6914f78a48fe0f86b838",
            "linux_arm64": "13638d7b03ca2204ca4d19075d6f66c04731d84e131cace35ea2b59ec38a0143",
        },
    },
}

for name, metadata in tools.items():
    version = metadata["version"]
    for a, h in metadata["hashes"].items():
        native = f"{CONFIG.OS}_{CONFIG.ARCH}" == a
        remote_file(
            name = name if native else f"{name}_{a}",
            binary = True,
            hashes = [h],
            url = f"https://github.com/please-build/cc-rules/releases/download/{name}-v{version}/{name}-{version}-{a}",
            visibility = ["PUBLIC"] if native else None,
        )
