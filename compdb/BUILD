subinclude("//build_defs:cc")

cc_binary(
    name = "compdb",
    srcs = ["compdb.cc"],
    cflags = [
        "-Icompdb/subprocess",
        "-Icompdb/json/single_include",
    ],
    deps = [
        ":json",
        ":subprocess",
    ],
)

github_repo(
    name = "subprocess",
    build_file = "subprocess.build",
    hashes = ["4dd85ccdbf5c78dee8321dc3587c2c4970793abc394c4b9288debc8825b01571"],
    repo = "arun11299/cpp-subprocess",
    revision = "v2.0",
    strip_prefix = "cpp-subprocess-2.0",
)

github_repo(
    name = "json",
    build_file = "json.build",
    hashes = ["83258a1aa68921b1d82ca3fcfeef04f27471d33cc0e51d516109d161916223f2"],
    repo = "nlohmann/json",
    revision = "v3.10.4",
)
