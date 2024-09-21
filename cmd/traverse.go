package cmd

import (
	"fmt"
	"os"
	"path"
)

// Traverse recursively the specified root folder
// Gather the list of files and the list of subfolders
// If there are files to be processed it calls ProduceEmbedGo()
func TraverseDir(root string) error {
	files := make([]string, 0)
	folders := make([]string, 0)

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// First determine files from folders
	// Note: .go files are not considered to be embedded
	for _, entry := range entries {
		if entry.IsDir() {
			gof := path.Join(root, entry.Name(), "embed.go")
			if DryRun {
				fmt.Printf("Would have deleted file %s\n", gof)
			} else {
				os.Remove(gof)
			}
			if OkToEmbedDir(entry.Name()) {
				folders = append(folders, entry.Name())
			}
		} else {
			if OkToEmbedFile(entry.Name()) {
				files = append(files, path.Base(entry.Name()))
			}
		}
	}

	// If there are files, then produce the embed.go file
	if len(files) > 0 {
		err := ProduceEmbedGo(root, files)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}

	// Now process all the folders
	for _, folder := range folders {
		if OkToEmbedDir(folder) {
			err := TraverseDir(path.Join(root, folder))
			if err != nil {
				return fmt.Errorf("%s", err.Error())
			}
		}
	}
	return nil
}
