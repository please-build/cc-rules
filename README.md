C and C++ rules for Please
==========================

This repo defines a [Please](https://please.build) plugin for building C and C++ code.

C and C++ build actions are very similar and can be intermingled freely (following the
usual language rules); the main reason for separating them is to allow defining different
tools and compiler flags for the two.

There are several different targets available providing different rules:


//build_defs:cc
---------------

Contains C++ rules. These use `cpp_tool`, `default_opt_cppflags`, `default_dbg_cppflags`
and `test_main`.
The following rules are defined.

 - `cc_library`
 - `cc_binary`
 - `cc_test`
 - `cc_object`
 - `cc_static_library`
 - `cc_shared_object`
 - `cc_module` (N.B. this is still experimental)

See the docstring for each rule for more specific detail on what they each do.


//build_defs:c
--------------

Contains C rules. These use `cc_tool`, `default_opt_cflags` and `default_dbg_cflags`.
The following rules are defined.

 - `c_library`
 - `c_binary`
 - `c_test`
 - `c_object`
 - `c_static_library`
 - `c_shared_object`


//build_defs:cc_embed_binary
----------------------------

Contains rules for embedding files into an object file that can be linked into a
binary and loaded at runtime. There are both C and C++ variants that work similarly
to the other build actions for each language.

These have the extra config value `default_namespace` to set the default namespace to
generate; it can also be overridden per target.

When building on OSX, the config value `asm_tool` is also used (by default this is `nasm`).
On other platforms this is not required.

 - `c_embed_binary`
 - `cc_embed_binary`


General notes
-------------

These are very much based on GCC and Clang; while it would be theoretically possible
to support MSVC the flag structure would need to change fairly dramatically (and hence it
may be easier to support them as a totally parallel set of rules or even a different plugin).
In practice this would also require Windows support for Please generally.
