// Binary generate_test_data generates test data files for the TestParseOutput test.
//
// generate_test_data is driven by a toolspec file, a JSON-encoded file specifying the paths to the compilers and
// linkers whose output should be captured (as well as how please_cc is expected to identify them). generate_test_data
// then invokes each possible compiler-linker combination, as well as each of the linkers individually, and writes
// ".test_data" files describing the output from each invocation and the expected result of parsing it. The files are
// written to the test_data/ directory by default, although this can be overridden with the -o option.
//
// Lines beginning with "#" (excluding whitespace) in toolspec files are treated as comments and are stripped before
// the file is decoded.
package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"

	"github.com/thought-machine/go-flags"
)

var regexpComment = regexp.MustCompile(`(?m)^\s*#.*?$`)

var opts struct {
	OutputDir string `short:"o" long:"output_dir" value-name:"DIR" default:"test_data" description:"Write test data to DIR"`
	Args struct {
		ToolSpec string `positional-arg-name:"TOOLSPEC_FILE" description:"Read toolspec from TOOLSPEC_FILE" required:"true"`
	} `positional-args:"true"`
}

// toolspec describes the C/C++ compilers and linkers from which test cases should be generated.
//
// generate_test_data runs each compiler multiple times, configuring it to invoke a different linker each time; it then
// runs each linker independently of a compiler. The total number of test cases produced from a single toolspec is
// therefore (len(toolspec.Compilers) * len(toolspec.Linkers)) + len(toolspec.Linkers), assuming no linkers are excluded
// for any compilers (see toolspec.Compilers.Exclude).
var toolspec struct {
	// Compilers identifies the C/C++ compilers for which output should be collected for test cases. Keys are paths to
	// executable files that invoke compilers; they may or may not be wrappers for other executables.
	Compilers map[string]struct {
		// ID is a string, unique to this toolspec, that identifies this compiler. The string forms part of the names of the
		// test files involving this compiler that generate_test_data writes.
		ID string `json:"id"`

		// Name is the compiler name and version that Tool.String is expected to return after the compiler's output has been
		// parsed.
		Name string `json:"name"`

		// Exclude is a list of IDs of linkers for which test cases should not be generated in combination with this compiler.
		//
		// generate_test_data combines linkers with compilers by creating an "ld" symlink that points to the linker in a
		// temporary directory, then prepending the path to the temporary directory to the compiler's search path. This is the
		// only universally-supported method of invoking an arbitrary linker among the compilers that Tool supports, but it
		// relies on the compiler being configured to use the relative path "ld" as its default compiler at build time - if
		// this is not the case, the compiler may invoke a different linker than the one generate_test_data expects. Exclude
		// can be defined in these circumstances to avoid generating bogus test cases.
		Exclude []string `json:"exclude"`
	} `json:"compilers"`

	// Linkers identifies the linkers for which output should be collected for test cases. Keys are paths to executable
	// files that invoke linkers; they may or may not be wrappers for other executables.
	Linkers map[string]struct {
		// ID is a string, unique to this toolspec, that identifies this linker. The string forms part of the names of the
		// test files involving this linker that generate_test_data writes.
		ID string `json:"id"`

		// Name is the linker name and version that Tool.String is expected to return after the linker's output has been
		// parsed.
		Name string `json:"name"`
	} `json:"linkers"`
}

func main() {
	_, err := flags.NewParser(&opts, flags.Default - flags.PrintErrors).Parse()
	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			fmt.Printf("%v", err)
			os.Exit(0)
		}
		log.Fatalf("Failed to parse command line arguments: %v", err)
	}

	toolspecBytes, err := os.ReadFile(opts.Args.ToolSpec)
	if err != nil {
		log.Fatalf("Failed to read toolspec file: %v", err)
	}
	err = json.Unmarshal(filterComments(toolspecBytes), &toolspec)
	if err != nil {
		log.Fatalf("Failed to parse toolspec file: %v", err)
	}

	for linkerPath, linkerData := range toolspec.Linkers {
		for compilerPath, compilerData := range toolspec.Compilers {
			if compilerData.Exclude != nil && slices.ContainsFunc(compilerData.Exclude, func(ex string) bool {
				matched, _ := filepath.Match(ex, linkerData.ID)
				return matched
			}) {
				continue
			}

			if err := run(
				compilerPath,
				linkerPath,
				fmt.Sprintf("cc=%q ld=%q", compilerData.Name, linkerData.Name),
				filepath.Join(opts.OutputDir, compilerData.ID+"-"+linkerData.ID+".test_data"),
			); err != nil {
				log.Fatalf("Failed to generate test data for %s / %s: %v", compilerData.ID, linkerData.ID, err)
			}
		}

		if err := run(
			"",
			linkerPath,
			fmt.Sprintf("ld=%q", linkerData.Name),
			filepath.Join(opts.OutputDir, linkerData.ID+".test_data"),
		); err != nil {
			log.Fatalf("Failed to generate test data for %s: %v", linkerData.ID, err)
		}
	}
}

// filterComments removes lines prepended with `#` from a byte slice.
func filterComments(in []byte) []byte {
	return regexpComment.ReplaceAll(in, []byte{})
}

// run invokes the linker at the given path - indirectly via the compiler, if a compiler path is also specified - and
// writes the output on the standard streams to the test file at testPath, overwriting any existing file. The compiler
// and/or linker are invoked with options that cause them to identify themselves on the standard streams. expected
// should be a string in one of the following forms:
//
// - `cc="CC_TOOL_STRING" ld="LD_TOOL_STRING"` (if compiler is non-nil)
// - `ld="LD_TOOL_STRING"` (if compiler is nil)
//
// where `CC_TOOL_STRING` is the compiler's identity as it would be expected to be reported by `Tool.String`, and
// `LD_TOOL_STRING` is the linker's identity as it would be expected to be reported by `Tool.String`.
func run(compiler, linker, expected, testPath string) error {
	var (
		tool string
		args []string
	)
	if compiler != "" {
		tool = compiler
		linkerDir, err := os.MkdirTemp("", "ld.")
		if err != nil {
			return err
		}
		defer os.RemoveAll(linkerDir)
		ldLink := filepath.Join(linkerDir, "ld")
		if err := os.Symlink(linker, ldLink); err != nil {
			return fmt.Errorf("symlink %s to %s: %v", linker, ldLink, err)
		}
		// -B prepends the temporary directory to the compiler's search path, so it finds the linker at "ld" in that
		// directory. -v causes the compiler to identify itself on one of the standard streams; -Wl,-v passes -v through to
		// the linker, causing it to do the same.
		args = []string{"-B", linkerDir, "-v", "-Wl,-v"}
	} else {
		tool = linker
		// We're invoking the linter directly, so there's no need for the complexity above here - just ask it to identify
		// itself on one of the standard streams.
		args = []string{"-v"}
	}

	var toolOut, toolErr bytes.Buffer
	cmd := exec.Command(tool, args...)
	cmd.Stdout = &toolOut
	cmd.Stderr = &toolErr
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	exp, err := os.OpenFile(testPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %v", testPath, err)
	}
	if _, err := exp.WriteString(fmt.Sprintf("#expected: %s\n#stdout:\n", expected)); err != nil {
		return fmt.Errorf("write %s: %v", testPath, err)
	}
	if _, err := exp.Write(toolOut.Bytes()); err != nil {
		return fmt.Errorf("write %s: %v", testPath, err)
	}
	if _, err := exp.WriteString("\n#stderr:\n"); err != nil {
		return fmt.Errorf("write %s: %v", testPath, err)
	}
	if _, err := exp.Write(toolErr.Bytes()); err != nil {
		return fmt.Errorf("write %s: %v", testPath, err)
	}
	return nil
}
