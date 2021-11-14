cc_binary(
    name = "compdb",
    srcs = ["compdb.cc"],
    cflags = ["-Icompdb/subprocess"],
    deps = [
        ":subprocess",
    ],
)

github_repo(
    name = "subprocess",
    repo = "arun11299/cpp-subprocess",
    revision = "v2.0",
    hashes = ["4dd85ccdbf5c78dee8321dc3587c2c4970793abc394c4b9288debc8825b01571"],
    strip_prefix = "cpp-subprocess-2.0",
    build_file = "subprocess.build",
)
