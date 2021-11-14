JSON Compilation Database
=========================

Contains a small program to create a JSON compilation database file for the current project
which can be read by Clang / libclang etc.

See https://clang.llvm.org/docs/JSONCompilationDatabase.html for more information about
the format.


Limitations
-----------

At present it does not build any dependencies so the some completions may not
work if they have dynamically generated inputs etc (similarly they will not
work after a `plz clean` which we can't do much about). This will be fixed after a
`plz build` of the relevant actions, after which all the dependencies will be present.

If you're using a [remote execution](https://please.build/remote_builds.html) profile
you may need to do `plz build --download` to enforce downloading particular outputs.
