package file_test

import (
	"os"
	"testing"

	"github.com/pieni-2-organiser/internal/file"
	"github.com/pieni-2-organiser/internal/fs"
	"github.com/stretchr/testify/assert"
)

func TestCleanupSys(t *testing.T) {

	type testCases struct {
		input    *fs.Contents
		expected *fs.Contents
	}

	cases := map[string]testCases{
		"empty": {
			input:    &fs.Contents{},
			expected: &fs.Contents{},
		},
		"no sys files": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
			expected: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
		},
		"with thumbs.db": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt", "Thumbs.db"},
			},
			expected: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
		},
		"with .DS_Store": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt", ".DS_Store"},
			},
			expected: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
		},
		"with both sys files": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt", "Thumbs.db", ".DS_Store"},
			},
			expected: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
		},
		"with mixed case": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt", "Thumbs.db", "thumbs.db", ".DS_Store", ".ds_store"},
			},
			expected: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt"},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := file.CleanupSys(tc.input)
			assert.Len(t, result.Files, len(tc.expected.Files), "expected number of files to match")
			assert.ElementsMatch(t, result.Files, tc.expected.Files, "expected files to match")
		})
	}
}

func TestNewRenameMap(t *testing.T) {
	type testCases struct {
		input    *fs.Contents
		expected file.RenameMap
	}

	cases := map[string]testCases{
		"empty": {
			input:    &fs.Contents{},
			expected: file.RenameMap{},
		},
		"single file": {
			input: &fs.Contents{
				Files: []string{"file1.txt"},
			},
			expected: file.RenameMap{"file1.txt": "00000000.txt"},
		},
		"multiple files same extension": {
			input: &fs.Contents{
				Files: []string{"file1.txt", "file2.txt", "file3.txt"},
			},
			expected: file.RenameMap{
				"file1.txt": "00000000.txt",
				"file2.txt": "00000001.txt",
				"file3.txt": "00000002.txt",
			},
		},
		"multiple files different extensions": {
			input: &fs.Contents{
				Files: []string{"image1.jpg", "image2.jpg", "document.pdf"},
			},
			expected: file.RenameMap{
				"image1.jpg":   "00000000.jpg",
				"image2.jpg":   "00000001.jpg",
				"document.pdf": "00000000.pdf",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := file.NewRenameMap(tc.input)
			assert.Len(t, result, len(tc.expected), "expected rename map length to match")
			assert.Equal(t, tc.expected, result, "expected rename map to match")
		})
	}
}

func TestCopyWithRename(t *testing.T) {
	type input struct {
		RenameMap file.RenameMap
		Contents  fs.Contents
	}
	type testCases struct {
		input    input
		expected file.RenameMap
	}

	// Define the expected rename maps for the test cases
	// this is to make modification of these maps easier since they're used in 2 places
	singleFileMap := file.RenameMap{"fixtures/parent01/file1.txt": "00000000.txt"}
	multifileMap := file.RenameMap{
		"fixtures/parent01/file1.txt": "00000000.txt",
		"fixtures/parent02/file2.txt": "00000001.txt",
		"fixtures/file3.txt":          "00000002.txt",
	}

	cases := map[string]testCases{
		"empty": {
			input: input{
				RenameMap: file.RenameMap{},
				Contents:  fs.Contents{},
			},
			expected: file.RenameMap{},
		},
		"single_file": {
			input: input{
				RenameMap: singleFileMap,
				Contents: fs.Contents{
					Files: []string{"fixtures/parent01/file1.txt"},
				},
			},
			expected: singleFileMap,
		},
		"multiple_files_same_extension": {
			input: input{
				RenameMap: multifileMap,
				Contents: fs.Contents{
					Files: []string{
						"fixtures/parent01/file1.txt",
						"fixtures/parent02/file2.txt",
						"fixtures/file3.txt",
					},
				},
			}, expected: multifileMap,
		}}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			workDir := "./fixtures/_workDir_" + name

			// cleanup working directory before each test iteration
			os.RemoveAll(workDir)
			// cleanup working directory when the test iteration is done
			defer os.RemoveAll(workDir)

			// create working directory for the test iteration
			err := os.Mkdir(workDir, 0755)
			assert.NoError(t, err, "expected no error in setting up the test working directory")

			result := file.CopyWithRename(workDir, &tc.input.Contents, tc.input.RenameMap)
			assert.NoError(t, result, "expected no error during copy with rename")

			// assert content of the working directory matches the expected rename map
			for original, renamed := range tc.input.RenameMap {
				// construct the expected path in the working directory
				expectedPath := workDir + "/" + renamed

				// check if the file exists in the working directory
				_, err := os.Stat(expectedPath)
				assert.NoError(t, err, "expected file %s to exist in %s", renamed, workDir)

				// check if the file content matches the original file
				originalContent, err := os.ReadFile(original)
				assert.NoError(t, err, "expected no error reading original file %s", original)

				renamedContent, err := os.ReadFile(expectedPath)
				assert.NoError(t, err, "expected no error reading renamed file %s", expectedPath)

				assert.Equal(t, originalContent, renamedContent, "expected content of %s to match %s", expectedPath, original)
			}
		})
	}
}
