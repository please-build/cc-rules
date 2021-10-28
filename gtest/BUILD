github_repo(
    name = "gtest",
    bazel_compat = True,
    repo = "google/googletest",
    revision = "release-1.8.1",
)

cc_library(
    name = "main",
    deps = [
        "///gtest/gtest//:gtest_main",
    ],
    visibility = ["PUBLIC"],
)
