github_repo(
    name = "unittest_cpp",
    repo = "unittest-cpp/unittest-cpp",
    hashes = ["3ba4f8b6c6e75ca8c79dd6008c78e95fb08554fbf2e0f1195fb78a4bf5a8a805"],
    build_file = "unittest.build",
    revision = "v2.0.0",
    strip_prefix = "unittest-cpp-2.0.0",
)

cc_library(
    name = "main",
    srcs = ["unittest_main.cc"],
    visibility = ["PUBLIC"],
    deps = ["///unittest-pp/unittest_cpp//:unittest_cpp"],
)
