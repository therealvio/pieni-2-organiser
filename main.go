package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type DirContents struct {
	Files []string
	Dirs  []string //TODO: dunno if we need dirs yet
}

// map of old file names to new file names
type RenameMap map[string]string

func main() {
	flag.Parse()
	root := flag.Arg(0)

	originalDirContents, err := DirectoryContents(root)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	CleanupSysFiles(originalDirContents)

	renameMap := NewRenameMap(originalDirContents)

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	workDir := cwd + "/pieni-2-organiser-output"
	err = CreateWorkDir(workDir)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = CopyFiles(workDir, originalDirContents, renameMap)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func DirectoryContents(root string) (*DirContents, error) {
	dirContents := &DirContents{}

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

// remove system files like .DS_Store and Thumbs.db from the dirContents struct, we're not interested in managing these
func CleanupSysFiles(dirContents *DirContents) *DirContents {
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

func NewRenameMap(directoryContents *DirContents) RenameMap {
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

func CopyFiles(workDir string, directoryContents *DirContents, renameMap RenameMap) error {
	// copy files into a temporary directory
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
