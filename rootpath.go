package xpress

import (
	"os"
	"path/filepath"
)

func GetRootPath() string {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Find the project's root directory by going up the directory tree
	for {
		// Check if a certain file or directory exists in the current working directory
		_, err := os.Stat(filepath.Join(wd, "go.mod"))
		if err == nil {
			// The "go.mod" file exists, which indicates the root of a Go module.
			break
		}

		// Move up one directory level
		wd = filepath.Dir(wd)

		// Check if we've reached the filesystem root (e.g., on Windows, it would be "C:\")
		if wd == "/" || wd == "\\" {
			// We've reached the root, and "go.mod" was not found.
			// Handle this case as needed.
			break
		}
	}

	// The variable 'wd' now contains the project's root path.
	return wd
}
