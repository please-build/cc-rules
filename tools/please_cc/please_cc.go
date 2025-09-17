// Binary please_cc implements a tool for invoking C/C++ compilers and linkers using command line options that are
// dynamically selected depending on the compiler or linker being invoked.
//
// please_cc operates in one of two modes, indicated by the value of the first argument passed to it:
//
// - `cc` ("compiler mode") invokes the compiler - and, implicitly, the linker - in a C or C++ toolchain.
// - `ld` ("linker mode") invokes a linker independently of a C or C++ toolchain.
//
// The second argument identifies the (relative or absolute) path to the compiler or linker. If the path is relative,
// please_cc follows the usual conventions for converting relative paths into absolute paths before invoking the
// compiler or linker.
//
// The remaining arguments are command line options to be passed through to the compiler or linker. Arguments surrounded
// by `{{ ` and ` }}` are evaluated as expressions (see [github.com/please-build/cc-rules/tools/please_cc/expr]) after
// those delimiters have been stripped; the result of each evaluation is substituted into the list of command line
// options passed to the compiler or linker. The identity of the compiler and/or linker is described in the environment
// in which the expression is evaluated, which allows for different command line options to be passed through for
// different compilers and/or linkers.
//
// For example:
//
// - Only make warnings fatal when compiling with GCC:
//   `please_cc cc cctool -o example '{{ gcc ? "-Werror" }}' example.c`
// - Compile C++ modules with GCC and Clang <= 15 with the `-fmodules-ts` option, or Clang >= 16 with the (equivalent,
//   non-deprecated) `-std=c++20` option:
//   `please_cc cc c++tool -o example '{{ gcc || (clang && clang <= 15) ? "-fmodules-ts" : "-std=c++20" }}' example.cc`
// - Link objects using the old ld64 code path if linking with Apple's new ld linker (enabled using `-ld64` prior to
//   Xcode 15.1, and `-ld_classic` from Xcode 15.1 onwards):
//   `please_cc ld ldtool '{{ appleld ? (appleld >= 1022.1 ? "-ld_classic" : "-ld64") }}' obj1.o obj2.o -o example`
//
// please_cc is known to be compatible with the following compilers and linkers:
//
// | Tool        | Minimum supported version | Expression language identifier |
// | ----------- | ------------------------- | ------------------------------ |
// | GCC         | 9                         | `gcc`                          |
// | Clang       | 11                        | `clang`                        |
// | Apple Clang | 12                        | `aclang`                       |
// | GNU ld      | 2.38                      | `gnuld`                        |
// | GNU gold    | 1.15                      | `gold`                         |
// | LLD         | 11                        | `lld`                          |
// | ld64        | 609.8                     | `ld64`                         |
// | Apple ld    | 1015.7                    | `appleld`                      |
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/please-build/cc-rules/tools/please_cc/cctool"
	"github.com/please-build/cc-rules/tools/please_cc/expr"
)

// exprPrefix and exprSuffix are the string prefix and suffix that indicate that a given command line argument is to be
// evaluated as an expression using please_cc's expression language.
const (
	exprPrefix = "{{ "
	exprSuffix = " }}"
)

func main() {
	toolType, tool, args, err := parseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var (
		realArgs = make([]string, 0)
		exprEnv map[string]any
	)
	for _, arg := range args {
		// We only need to go through the expense of identifying the compiler and/or linker if any of the command line
		// arguments we're passing to the tool are expressions - as soon as we see one, identify the tool(s) and template the
		// arguments.
		if strings.HasPrefix(arg, exprPrefix) && strings.HasSuffix(arg, exprSuffix) {
			if exprEnv == nil {
				var (
					compiler, linker *cctool.Tool
					err              error
				)
				switch toolType {
				case cctool.Compiler:
					compiler, linker, err = cctool.IdentifyCompiler(tool, args)
				case cctool.Linker:
					linker, err = cctool.IdentifyLinker(tool)
				}
				if err != nil {
					log.Fatalf("%s: %v\n", toolType, err)
				}
				if compiler != nil {
					log.Printf("Identified C/C++ compiler as %s\n", compiler)
				}
				log.Printf("Identified linker as %s\n", linker)
				exprEnv = environment(compiler, linker)
			}
			exprArg := strings.TrimSuffix(strings.TrimPrefix(arg, exprPrefix), exprSuffix)
			exprArgs, err := expr.Evaluate(exprArg, exprEnv)
			if err != nil {
				log.Fatalf("Failed to evaluate expression `%s`: %v\n", exprArg, err)
			}
			realArgs = append(realArgs, exprArgs...)
		} else {
			realArgs = append(realArgs, arg)
		}
	}
	// Re-invoke the tool with the templated command line arguments. please_cc doesn't need to wait on the child process
	// (the caller will handle failures if there are any), so we can just invoke it with (a Go equivalent to) execvp().
	log.Printf("Running `%s %s`\n", tool, strings.Join(realArgs, " "))
	if err := execvp(tool, realArgs); err != nil {
		log.Fatalf("Failed to run tool: %v\n", err)
	}
}

func parseArgs(args []string) (cctool.Type, string, []string, error) {
	if len(args) < 3 {
		return "", "", nil, fmt.Errorf("Usage: %s [cc|ld] [TOOL_PATH] [TOOL_ARGS]...", args[0])
	}
	toolArgs := []string{}
	if len(args) > 3 {
		toolArgs = args[3:]
	}
	switch cctool.Type(args[1]) {
	case cctool.Compiler:
		return cctool.Compiler, args[2], toolArgs, nil
	case cctool.Linker:
		return cctool.Linker, args[2], toolArgs, nil
	default:
		return "", "", nil, fmt.Errorf("Usage: %s [cc|ld] [TOOL_PATH] [TOOL_ARGS]...", args[0])
	}
}

func environment(compiler, linker *cctool.Tool) map[string]any {
	env := make(map[string]any)
	if compiler != nil {
		switch compiler.Name {
		case cctool.GCC:
			env["gcc"] = compiler.Version
		case cctool.AppleClang:
			env["aclang"] = compiler.Version
		case cctool.Clang:
			env["clang"] = compiler.Version
		}
	}
	switch linker.Name {
	case cctool.GNUld:
		env["gnuld"] = linker.Version
	case cctool.LLD:
		env["lld"] = linker.Version
	case cctool.AppleLd:
		env["appleld"] = linker.Version
	case cctool.Ld64:
		env["ld64"] = linker.Version
	case cctool.Gold:
		env["gold"] = linker.Version
	}
	return env
}

// execvp mimics libc's execvp() function - it duplicates (or, at least, approximates) the behaviour of the shell in
// searching PATH for an executable with the given file name and invokes it with the given list of arguments. If this is
// successful, it never returns.
func execvp(file string, args []string) error {
	execFile, err := exec.LookPath(file)
	if err == nil {
		execFile, err = filepath.Abs(execFile)
	}
	if err != nil {
		return err
	}
	return syscall.Exec(execFile, append([]string{file}, args...), os.Environ())
}
