package(default_visibility = ["PUBLIC"])

subinclude(
    "///go//build_defs:go",
    "//build_defs:arch",
)

go_toolchain(
    name = "toolchain",
    install_std = False,
    architectures = SUPPORTED_ARCHITECTURES,
    version = "1.24.6",
)

go_stdlib(
    name = "std",
    tags = [
        "netgo",
        "osusergo",
    ],
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/PaesslerAG/gval",
    version = "v1.2.4",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/PaesslerAG/jsonpath",
    version = "v0.1.0",
)

go_repo(
    licences = ["MIT"],
    module = "github.com/shopspring/decimal",
    version = "v1.4.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/generikvault/gvalstrings",
    version = "v0.0.0-20180926130504-471f38f0112a",
)

go_repo(
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.11.1",
)

go_repo(
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.2-0.20180830191138-d8f796af33cc",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.1-0.20181226105442-5d4384ee4fb2",
)

go_repo(
    licences = ["MIT"],
    module = "github.com/stretchr/objx",
    version = "v0.5.2",
)

go_repo(
    licences = [
        "MIT",
        "Apache-2.0",
    ],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.1",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.7.0",
)

go_repo(
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sys",
    version = "v0.36.0",
)
