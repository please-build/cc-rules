package cctool

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDataPath = "tools/please_cc/cctool/test_data"

var regexpTestData = regexp.MustCompile(`(?s)^#expected:(?: cc="(?P<cc>.+?)")? ld="(?P<ld>.+?)"\n#stdout:\n(?P<stdout>.*?)\n#stderr:\n(?P<stderr>.*)$`)

func TestFilterLinkerArgs(t *testing.T) {
	for desc, test := range map[string]struct {
		Args     []string
		Expected []string
	}{
		"Empty list": {
			Args:     []string{},
			Expected: []string{},
		},
		"-B in single option": {
			Args:     []string{"-Wall", "-Bbin", "-Werror"},
			Expected: []string{"-Bbin"},
		},
		"-B over two options": {
			Args:     []string{"-Wall", "-B", "bin", "-Werror"},
			Expected: []string{"-B", "bin"},
		},
		"-fuse-ld with linker name": {
			Args:     []string{"-Wall", "-fuse-ld=gold", "-Werror"},
			Expected: []string{"-fuse-ld=gold"},
		},
		"-fuse-ld with relative path": {
			Args:     []string{"-Wall", "-fuse-ld=ld_dir/ld", "-Werror"},
			Expected: []string{"-fuse-ld=ld_dir/ld"},
		},
		"-fuse-ld with absolute path": {
			Args:     []string{"-Wall", "-fuse-ld=/usr/bin/ld.lld", "-Werror"},
			Expected: []string{"-fuse-ld=/usr/bin/ld.lld"},
		},
		"Invalid -fuse-ld (no =)": {
			Args:     []string{"-Wall", "-fuse-ld", "gold", "-Werror"},
			Expected: []string{},
		},
		"Invalid -fuse-ld (two leading -s)": {
			Args:     []string{"-Wall", "--fuse-ld=gold", "-Werror"},
			Expected: []string{},
		},
		"--ld-path with relative path": {
			Args:     []string{"-Wall", "--ld-path=ld_dir/ld.gold", "-Werror"},
			Expected: []string{"--ld-path=ld_dir/ld.gold"},
		},
		"--ld-path with absolute path": {
			Args:     []string{"-Wall", "--ld-path=/usr/bin/ld.gold", "-Werror"},
			Expected: []string{"--ld-path=/usr/bin/ld.gold"},
		},
		"Invalid --ld-path (no =)": {
			Args:     []string{"-Wall", "--ld-path", "/usr/bin/ld.gold", "-Werror"},
			Expected: []string{},
		},
		"Invalid --ld-path (only one leading -)": {
			Args:     []string{"-Wall", "-ld-path=/usr/bin/ld.gold", "-Werror"},
			Expected: []string{},
		},
	} {
		assert.EqualValues(t, test.Expected, filterLinkerArgs(test.Args), desc)
	}
}

func TestParseOutput(t *testing.T) {
	err := filepath.WalkDir(testDataPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.Type().IsRegular() || !strings.HasSuffix(path, ".test_data") {
			return nil
		}
		relPath := strings.TrimPrefix(path, testDataPath+"/")
		t.Run(strings.TrimSuffix(relPath, ".test_data"), func(t *testing.T) {
			testData, err := os.ReadFile(path)
			assert.NoErrorf(t, err, "Failed to read test file: %s", path)
			test := regexpTestData.FindSubmatch(testData)
			assert.NotNilf(t, test, "Malformed test file: %s", path)
			expectedCompiler := string(test[regexpTestData.SubexpIndex("cc")])
			expectedLinker := string(test[regexpTestData.SubexpIndex("ld")])
			stdout := test[regexpTestData.SubexpIndex("stdout")]
			stderr := test[regexpTestData.SubexpIndex("stderr")]
			actualCompiler, actualLinker := parseOutput(stdout, stderr)
			if expectedCompiler == "" {
				assert.Nil(t, actualCompiler, "Compiler is nil")
			} else {
				assert.NotNil(t, actualCompiler, "Compiler is not nil")
				assert.Equal(t, expectedCompiler, actualCompiler.String(), "Compiler correctly identified")
			}
			assert.NotNil(t, actualLinker, "Linker is not nil")
			assert.Equal(t, expectedLinker, actualLinker.String(), "Linker correctly identified")
		})
		return nil
	})
	assert.NoError(t, err, testDataPath+" walked successfully")
}
