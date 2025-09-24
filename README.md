# C and C++ rules for Please

This repo defines a [Please](https://please.build) plugin for building C and C++ code.

C and C++ build actions are very similar and can be intermingled freely (following the
usual language rules); the main reason for separating them is to allow defining different
tools and compiler flags for the two.


## Basic usage


```python
# Add this to plugins/BUILD 
plugin_repo(
    name = "cc",
    revision = "v1.0.0",
)

# Subinclude from the plugin
subinclude("///plugins/cc//build_defs:cc")

# Compile some C code
cc_library(
    name = "foo",
    srcs = ["foo.c"],
    # Headers are specified separately so they can be exposed to other rules
    hdrs = ["foo.h"],
)

cc_binary(
   name = "main",
   srcs = ["main.c"],
   deps = [":foo"],
)
```

## Build definitions

There are several different targets available providing different rules:


### //build_defs:cc

Contains the following C++ rules that use `cpp_tool`, `default_opt_cppflags`, `default_dbg_cppflags`
and `test_main`.

 - `cc_library()` 
 - `cc_binary()`
 - `cc_test()`
 - `cc_object()`
 - `cc_static_library()`
 - `cc_shared_object()`
 - `cc_module()` (N.B. this is still experimental)

And the following C rules that use `cc_tool`, `default_opt_cflags` and `default_dbg_cflags`:

- `c_library()`
- `c_binary()`
- `c_test()`
- `c_object()`
- `c_static_library()`
- `c_shared_object()`

See the docstring for each rule for more specific detail on what they each do.


## Configuration

This plugin can be configured by adding fields to the `[Plugin "cc"]` section in your 
`.plzconfig`. The available configuration settings are documented here.

### CCTool
The tool used by `c_xxx()` build definitions to compile C code. Defaults to `gcc`. 
```ini
[Plugin "cc"]
CCTool = clang
```

### CPPTool
The tool used by `cc_xxx()` build definitions to compile C++ code. Defaults to `g++`.
```ini
[Plugin "cc"]
CPPTool = clang++
```

### ARTool
The tool used to manipulate `.a` archives. Defaults to `ar`. 
```ini
[Plugin "cc"]
ARTool = ar
```

### DefaultOptCFlags
Default flags used to compile C code. Defaults to `-std=c99 -O3 -pipe -DNDEBUG -Wall -Werror`.
```ini
[Plugin "cc"]
DefaultOptCFlags = -std=c99
DefaultOptCFlags = -O3
```

### DefaultDbgCFlags 
Default flags used to compile C code for debugging. Defaults to `-std=c99 -g3 -pipe -DDEBUG -Wall -Werror`.
```ini
[Plugin "cc"]
DefaultDbgCFlags = -std=c99
DefaultDbgCFlags = -O3
```

### DefaultOptCppFlags
Default flags used to compile C++ code. Defaults to `-std=c++11 -O3 -pipe -DNDEBUG -Wall -Werror`.
```ini
[Plugin "cc"]
DefaultOptCppFlags = -std=c99
DefaultOptCppFlags =-O3
```

### DefaultDbgCppFlags 
Default flags used to compile C++ code for debugging. Defaults to `-std=c++11 -g3 -pipe -DDEBUG -Wall -Werror`.
```ini
[Plugin "cc"]
DefaultDbgCppFlags = -std=c99
DefaultDbgCppFlags = -O3
```

### DefaultLDFlags
Default flags to pass when linking C and C++ code. Defaults to `-lpthread -ldl`.
```ini
[Plugin "cc"]
DefaultLDFlags = -ldl
```

### PkgConfigPath
Controls the `PKG_CONFIG_PATH` environment variable used by `pkg_config`. Not set by default. 
```ini
[Plugin "cc"]
PackageConfigPath = /opt/toolchain/pkg_configs
```

### TestMain
A `cc_library()`, `c_library()`, or otherwise compatible rule containing the entry point to run tests. 
Defaults to `//unitest-pp:main` in this plugin. 

```ini
[Plugin "cc"]
TestMain = //third_party/cc:gtest_main
```

### DsymTool
On `macOS`, the tool used to create debug symbols. Defaults to `dsymutil`. 

```ini
[Plugin "cc"]
DsymTool = dsymutil
```

### DefaultNamespace
The default C++ namespace to use. By default, no namespace is used. 
```ini
[Plugin "cc"]
DefaultNamespace = foo
```

## General notes

These are very much based on GCC and Clang; while it would be theoretically possible
to support MSVC the flag structure would need to change fairly dramatically (and hence it
may be easier to support them as a totally parallel set of rules or even a different plugin).
In practice this would also require Windows support for Please generally.
