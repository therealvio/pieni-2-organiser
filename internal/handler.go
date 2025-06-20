package internal

import (
	"fmt"
	"os"

	"github.com/pieni-2-organiser/internal/file"
	"github.com/pieni-2-organiser/internal/fs"
)

func Handler(root string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("retrieving current working directory: %w", err)
	}
	workDir := cwd + "/pieni-2-organiser-output"

	originalDirContents, err := fs.GetContents(root)
	if err != nil {
		return fmt.Errorf("retrieving directory contents in %v: %w", workDir, err)
	}
	file.CleanupSys(originalDirContents)

	renameMap := file.NewRenameMap(originalDirContents)

	err = fs.CreateWorkDir(workDir)
	if err != nil {
		return fmt.Errorf("creating work directory: %w", err)
	}

	err = file.CopyWithRename(workDir, originalDirContents, renameMap)
	if err != nil {
		return fmt.Errorf("copying files with rename: %w", err)
	}

	return nil
}
