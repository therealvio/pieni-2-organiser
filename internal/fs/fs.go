package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type Contents struct {
	Files []string
	Dirs  []string //TODO: dunno if we need dirs yet
}

func GetContents(root string) (*Contents, error) {
	dirContents := &Contents{}

	// use a closure to enumerate files and directories and keep dirContents local
	enumerate := func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			dirContents.Dirs = append(dirContents.Dirs, path)

		} else {
			dirContents.Files = append(dirContents.Files, path)
		}

		return nil

	}

	err := filepath.WalkDir(root, enumerate)

	if err != nil {
		return nil, err
	}

	return dirContents, nil
}

func CreateWorkDir(workDir string) error {
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		err := os.Mkdir(workDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}

	return nil
}
