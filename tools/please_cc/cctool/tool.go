package cctool

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Type is a family of tools that this package can identify.
type Type string

const (
	// Compiler is a tool for compiling C or C++ code.
	Compiler Type = "cc"

	// Linker is a tool for linking libraries and objects into an executable.
	Linker = "ld"
)

// Name is the (unique) name of a C/C++ compiler or linker.
type Name string

const (
	// GCC is the [GNU Compiler Collection].
	//
	// [GNU Compiler Collection]: https://gcc.gnu.org
	GCC Name = "GCC"

	// Clang is the [C language family frontend for LLVM].
	//
	// [C language family frontend for LLVM]: https://clang.llvm.org
	Clang = "Clang"

	// AppleClang is [Apple's build of Clang], distributed as part of Xcode.
	//
	// [Apple's build of Clang]: https://developer.apple.com/xcode/cpp/
	AppleClang = "Apple Clang"

	// GNUld is the [GNU linker], distributed as part of GNU binutils.
	//
	// [GNU linker]: https://www.gnu.org/software/binutils/
	GNUld = "GNU ld"

	// Gold is [GNU gold], distributed as part of GNU binutils.
	//
	// [GNU gold]: https://www.gnu.org/software/binutils/
	Gold = "GNU gold"

	// LLD is the [LLVM Linker].
	//
	// [LLVM Linker]: https://lld.llvm.org
	LLD = "LLD"

	// Ld64 is Apple's [ld64] linker. In Xcode >= 15 it is referred to as ld-classic to distinguish it from Apple ld.
	//
	// [ld64]: https://github.com/apple-oss-distributions/ld64
	Ld64 = "ld64"

	// AppleLd is Apple's ld linker, introduced in Xcode 15. During the Xcode 15 beta cycle it was known as ld-prime.
	AppleLd = "Apple ld"
)

// Tool is a specific version of a C/C++ compiler or linker.
//
// The tools identifiable by please_cc all use dot-decimal notation for version numbers of stable releases. Development
// releases may contain other identifiers (such as Git commit hashes), but these are not captured in Tool; a specific
// instance of Tool may therefore not accurately distinguish proximate development builds of the same tool.
type Tool struct {
	// Name is the tool's name.
	Name Name

	// Version is the tool's version.
	Version Version
}

// String returns the tool's name and version number.
func (t *Tool) String() string {
	return fmt.Sprintf("%s %s", t.Name, t.Version)
}

// stream represents a standard stream on which a tool produces its output.
type stream int

const (
	stdout stream = iota
	stderr
)

// matcher couples a tool name/type with a regular expression that identifies that tool based on its output when invoked
// with the verbose (-v) command line option.
type matcher []struct {
	typ    Type
	name   Name
	regexp *regexp.Regexp
}

var matchers = map[stream]matcher{
	// These tools identify themselves via their stderr output:
	stderr: matcher{
		{ Compiler, GCC, regexp.MustCompile(`^gcc (?:version|\(GCC\)) (?P<version>[\d.]+)`) },
		// Apple Clang's identification line includes two version numbers: a semantic version number and an internal build
		// version number. The former is the one captured in the "version" group, on the basis that the detail of the latter
		// is unlikely to matter for please_cc's purposes.
		//
		// Apple Clang uses a different version numbering system to that of Clang builds derived from the open-source LLVM
		// project (see [Xcode toolchain versions]). Because the Clang matcher (below) also matches Apple Clang's
		// identification line, the Apple Clang matcher *must* be used before the Clang matcher.
		//
		// [Xcode toolchain versions]: https://en.wikipedia.org/wiki/Xcode#Toolchain_versions
		{ Compiler, AppleClang, regexp.MustCompile(`^Apple clang version (?P<version>[\d.]+)`) },
		// LLVM contains a configure-time option that enables vendors to prepend optional, arbitrary strings to the names of
		// its tools (e.g. [Compiler Explorer]'s Clang builds identify themselves simply as "clang", while the Clang builds in
		// Ubuntu's clang-* packages identify themselves as "Ubuntu clang"). This regular expression is vendor-agnostic.
		//
		// Note that, due to [regexp]'s lack of support for zero-length assertions, this regular expression also matches the
		// identification line outputted by Apple Clang; to accurately identify Apple Clang, this matcher must therefore be
		// used after the Apple Clang matcher.
		//
		// [Compiler Explorer]: https://godbolt.org
		{ Compiler, Clang, regexp.MustCompile(`^(?:[[:print:]]+ )?clang version (?P<version>[\d.]+)`) },
		// Between Xcode 15 and 15.2, Apple ld identified itself with the project name "dyld" (version 1015.7 in Xcode 15 and
		// 15.0.1, and version 1022.1 in Xcode 15.1 and 15.2). From Xcode 15.3 onwards, it identifies itself with the project
		// name "ld".
		//
		// When invoked with the `-ld64` (Xcode < 15.1) or `-ld_classic` (Xcode >= 15.1) options, Apple ld behaves like ld64,
		// although it continues to identify itself as Apple ld. Refer to the [Apple library primer] for more information.
		//
		// [Apple library primer]: https://developer.apple.com/forums/thread/715385
		{ Linker, AppleLd, regexp.MustCompile(`^@\(#\)PROGRAM:[[:print:]]+?\s+PROJECT:(?:dy)?ld-(?P<version>[\d.]+)$`) },
		{ Linker, Ld64, regexp.MustCompile(`^@\(#\)PROGRAM:[[:print:]]+?\s+PROJECT:ld64-(?P<version>[\d.]+)$`) },
	},
	// These tools identify themselves via their stdout output:
	stdout: matcher{
		// GNU binutils contains a configure-time option that enables vendors to append arbitrary strings to the names of its
		// tools (e.g. the GNU ld build in FreeBSD's devel/binutils package identifies itself simply as "GNU ld (GNU
		// Binutils)", while the ld build in Ubuntu's binutils package identifies itself as "GNU ld (GNU Binutils for
		// Ubuntu)"). This regular expression is vendor-agnostic.
		{ Linker, GNUld, regexp.MustCompile(`^GNU ld \(.*\) (?P<version>[\d.]+)$`) },
		// LLVM contains a configure-time option that enables vendors to prepend optional, arbitrary strings to the names of
		// its tools (e.g. the LLD builds in FreeBSD's devel/llvm* packages identify themselves simply as "LLD", while the LLD
		// builds in Ubuntu's lld-* packages identify themselves as "Ubuntu LLD"). This regular expression is vendor-agnostic.
		{ Linker, LLD, regexp.MustCompile(`^(?:[[:print:]]+ )?LLD (?P<version>[\d.]+).* \(compatible with GNU linkers\)$`) },
		// Although its source code is distributed as part of GNU binutils, gold is versioned independently of the other
		// tools. However, it identifies itself using both its own version number and the version number of the binutils
		// package in which it was distributed (e.g. gold 1.16 identifies itself as "GNU gold (GNU Binutils 2.38) 1.16" when
		// built from the binutils 2.38 source code, and as "GNU gold (GNU Binutils 2.44) 1.16" when built from the binutils
		// 2.44 source code). The gold version number is the one captured in the "version" group, on the basis that the
		// binutils version number is irrelevant for please_cc's purposes.
		//
		// GNU binutils contains a configure-time option that enables vendors to append arbitrary strings to the names of its
		// tools (e.g. the gold build in FreeBSD's devel/binutils package identifies itself simply as "GNU gold (GNU Binutils
		// [binutils version])", while the gold build in Ubuntu's binutils package identifies itself as "GNU gold (GNU
		// Binutils for Ubuntu [binutils version])"). This regular expression is vendor-agnostic.
		{ Linker, Gold, regexp.MustCompile(`^GNU gold \(.*\) (?P<version>[\d.]+)$`) },
	},
}

const reportURL = "https://github.com/please-build/cc-rules/issues/new"

// IdentifyCompiler invokes the C/C++ compiler at the given path and identifies both the name and version of the
// compiler and the name and version of the linker it is configured to use, taking into account the given command line
// options (which may influence the compiler's choice of linker). Both tools are identified based on their output to
// stdout and/or stderr.
func IdentifyCompiler(path string, args []string) (*Tool, *Tool, error) {
	// Every compiler and linker supported by cc-rules implements the -v option, which causes the tool to identify itself
	// either on stdout or stderr. -v causes the compiler to identify itself, while -Wl,-v passes -v through to the linker
	// and causes it to identify itself too. Although this is a nonsensical invocation of a C/C++ compiler toolchain (and
	// is therefore guaranteed to exit unsuccessfully, which means we have to ignore the process's exit code), this is the
	// simplest way to identify both tools in a single invocation.
	compilerArgs := append(filterLinkerArgs(args), "-v", "-Wl,-v")
	return identify(Compiler, path, compilerArgs)
}

// IdentifyLinker invokes the linker at the given path and identifies its name and version, based on its output to
// stdout and/or stderr.
func IdentifyLinker(path string) (*Tool, error) {
	_, linker, err := identify(Linker, path, []string{"-v"})
	return linker, err
}

func identify(typ Type, path string, args []string) (*Tool, *Tool, error) {
	var toolOut, toolErr bytes.Buffer
	tool := exec.Command(path, args...)
	tool.Stdout = &toolOut
	tool.Stderr = &toolErr
	if err := tool.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, nil, fmt.Errorf("execute %s: %w", path, err)
		}
	}
	compiler, linker := parseOutput(toolOut.Bytes(), toolErr.Bytes())
	if typ == Compiler && compiler == nil {
		return nil, nil, fmt.Errorf("failed to identify C/C++ compiler; please report the output of `%s %s` to %s", path, strings.Join(args, " "), reportURL)
	}
	if linker == nil {
		return nil, nil, fmt.Errorf("failed to identify linker; please report the output of `%s %s` to %s", path, strings.Join(args, " "), reportURL)
	}
	return compiler, linker, nil
}

// filterArgs returns the C/C++ compiler toolchain command line options in args that influence the toolchain's choice of
// linker.
//
// The C/C++ compilers supported by cc-rules honour a subset of the following [linker command line options], which are
// returned if present, regardless of whether the compiler in use actually supports them:
//
// - `-B[DIR]` or `-B [DIR]` (GCC, Clang): prepend `DIR` to compiler toolchain's `PATH`; if the compiler's default
//   linker is a relative path (it is typically, but not always, `ld`), `DIR` will implicitly be searched first.
// - `-fuse-ld=[NAME] (GCC, Clang): search `PATH` for `ld.NAME` and use the first result as the linker; `NAME` must be
//   one of `bfd`, `gold`, `lld` or `mold` for it to be accepted by current compilers, although filterLinkerArgs does
//   not enforce this.
// - `-fuse-ld=[PATH]` (Clang): use the linker at `PATH`; deprecated in favour of `--ld-path`.
// - `--ld-path=[PATH]` (Clang): use the linker at `PATH`.
//
// [linker command line options]: https://github.com/rust-lang/rust/issues/97402
func filterLinkerArgs(args []string) []string {
	linkerArgs := make([]string, 0)
	nextArg := false
	for _, arg := range args {
		if nextArg {
			linkerArgs = append(linkerArgs, arg)
			nextArg = false
		} else if arg == "-B" {
			linkerArgs = append(linkerArgs, arg)
			nextArg = true
		} else if strings.HasPrefix(arg, "-B") || strings.HasPrefix(arg, "-fuse-ld=") || strings.HasPrefix(arg, "--ld-path=") {
			linkerArgs = append(linkerArgs, arg)
		}
	}
	return linkerArgs
}

// parseOutput parses the output from a C/C++ compiler toolchain on stdout and stderr and identifies which compiler and
// linker produced the output by searching it for the first occurrence of a tool identification line for each tool type.
//
// Different tools identify themselves with output on different standard streams. parseOutput processes stdout first, as
// it contains only a single line identifying the linker in all but one case, which minimises the amount of processing
// of stderr that is required.
//
// The uses of MustParseVersion in this function are safe, since the only subexpressions matched by the tool
// identification line regular expressions are in dot-decimal notation.
func parseOutput(toolStdout, toolStderr []byte) (*Tool, *Tool) {
	var compiler, linker *Tool
	parseStream := func(str stream, out []byte) {
		for _, line := range bytes.Split(out, []byte("\n")) {
			for _, m := range matchers[str] {
				if m.typ == Compiler && compiler == nil {
					if sm := m.regexp.FindSubmatch(line); sm != nil {
						compiler = &Tool{
							Name:    m.name,
							Version: MustParseVersion(string(sm[m.regexp.SubexpIndex("version")])),
						}
					}
				}
			}
			for _, m := range matchers[str] {
				if m.typ == Linker && linker == nil {
					if sm := m.regexp.FindSubmatch(line); sm != nil {
						linker = &Tool{
							Name:    m.name,
							Version: MustParseVersion(string(sm[m.regexp.SubexpIndex("version")])),
						}
					}
				}
			}
			if compiler != nil && linker != nil {
				break
			}
		}
	}
	parseStream(stdout, toolStdout)
	if compiler != nil && linker != nil {
		return compiler, linker
	}
	parseStream(stderr, toolStderr)
	return compiler, linker
}
