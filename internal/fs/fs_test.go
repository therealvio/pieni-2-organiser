package fs_test

import (
	"testing"

	"github.com/pieni-2-organiser/internal/fs"
	"github.com/stretchr/testify/assert"
)

func TestGetContents(t *testing.T) {
	expectedResult := &fs.Contents{
		Files: []string{"../file/fixtures/file3.txt", "../file/fixtures/parent02/file2.txt", "../file/fixtures/parent01/file1.txt"},
	}

	dir := "../file/fixtures"
	result, err := fs.GetContents(dir)

	assert.NoError(t, err, "should not return an error when getting contents of a directory")
	assert.ElementsMatch(t, expectedResult.Files, result.Files, "should return the correct files in the directory")
}
