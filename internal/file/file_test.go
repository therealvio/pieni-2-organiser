package file_test

import (
	"testing"

	"github.com/pieni-2-organiser/internal/fs"
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
			{
				input: &fs.Contents{
					Files: []string{"file1.txt", "file2.txt"},
				},
				expected: &fs.Contents{},
			},
		},
		"with thumbs.db":      {},
		"with .DS_Store":      {},
		"with both sys files": {},
	}

}
