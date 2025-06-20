package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pieni-2-organiser/internal/fs"
)

type RenameMap map[string]string

// Remove system files like .DS_Store and Thumbs.db from the dirContents struct, we're not interested in managing these
func CleanupSys(dirContents *fs.Contents) *fs.Contents {
	// construct new slice of files - common go practice is to avoid modifying slices in place while iterating over them
	cleanFiles := make([]string, 0, len(dirContents.Files))
	for _, file := range dirContents.Files {
		base := filepath.Base(file)
		// if the file name is not a system file, add it to the cleanFiles slice
		if base != "Thumbs.db" && base != ".DS_Store" {
			cleanFiles = append(cleanFiles, file)
		}
	}
	dirContents.Files = cleanFiles

	return dirContents
}

func NewRenameMap(directoryContents *fs.Contents) RenameMap {
	renameMap := make(RenameMap)

	// Map extension to counter to support sequential renaming for different file types
	extCounters := make(map[string]int)

	files := directoryContents.Files
	sort.Strings(files)

	for _, file := range files {
		ext := filepath.Ext(file)
		counter := extCounters[ext]
		newBase := fmt.Sprintf("%08d%s", counter, ext)
		newName := filepath.Join(filepath.Dir(file), newBase)
		renameMap[file] = newName
		extCounters[ext] = counter + 1
	}

	return renameMap
}

// CopyWithRename files into the working directory
// This does not delete the original files, confirm the output is correct before deleting the original files
func CopyWithRename(workDir string, directoryContents *fs.Contents, renameMap RenameMap) error {
	for _, file := range directoryContents.Files {
		// get the new base name from the renameMap
		newBase := filepath.Base(renameMap[file])
		// construct the new path in the temp directory using the renamed base
		newPath := filepath.Join(workDir, newBase)

		// copy the file to the new path
		input, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}
		defer input.Close()

		output, err := os.Create(newPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", newPath, err)
		}
		defer output.Close()

		if _, err := output.ReadFrom(input); err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %w", file, newPath, err)
		}
	}

	return nil
}
