package replacelayer

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRoot(path string) (string, error) {
	root, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error finding directory: ", err.Error())
		return "", err
	}
	return root, nil
}

func GetFileNames(root string) map[string]bool {
	fileList := make(map[string]bool)

	files, err := os.ReadDir(root)

	if err != nil {
		fmt.Println("Error reading directory: ", root)
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		fileList[f.Name()] = true
	}

	return fileList
}
