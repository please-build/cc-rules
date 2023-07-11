subinclude("//build_defs:cc")

github_repo(
    name = "gtest",
    bazel_compat = True,
    repo = "google/googletest",
    revision = "release-1.8.1",
)

cc_library(
    name = "main",
    visibility = ["PUBLIC"],
    deps = [
        "///gtest/gtest//:gtest_main",
    ],
)
