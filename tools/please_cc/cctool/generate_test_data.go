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

var toolspec struct {
	Compilers map[string]struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Exclude []string `json:"exclude"`
	} `json:"compilers"`
	Linkers map[string]struct {
		ID   string `json:"id"`
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
				fmt.Sprintf("cc=\"%s\" ld=\"%s\"", compilerData.Name, linkerData.Name),
				filepath.Join(opts.OutputDir, compilerData.ID+"-"+linkerData.ID+".test_data"),
			); err != nil {
				log.Fatalf("Failed to generate test data for %s / %s: %v", compilerData.ID, linkerData.ID, err)
			}
		}

		if err := run(
			"",
			linkerPath,
			fmt.Sprintf("ld=\"%s\"", linkerData.Name),
			filepath.Join(opts.OutputDir, linkerData.ID+".test_data"),
		); err != nil {
			log.Fatalf("Failed to generate test data for %s: %v", linkerData.ID, err)
		}
	}
}

func filterComments(in []byte) []byte {
	return regexpComment.ReplaceAll(in, []byte{})
}

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
		args = []string{"-B", linkerDir, "-v", "-Wl,-v"}
	} else {
		tool = linker
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
