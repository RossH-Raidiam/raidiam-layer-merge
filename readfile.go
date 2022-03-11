package replacelayer

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRoot(path string) string {
	root, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error finding directory: ", err.Error())
		os.Exit(1)
	}
	return root
}

func GetFileNames(root string) (map[string]bool, error) {
	fileList := make(map[string]bool)

	files, err := os.ReadDir(root)

	if err != nil {
		fmt.Println("Error reading directory: ", root)
		fmt.Println("Error: ", err.Error())
		return nil, err
	}

	for _, f := range files {
		fileList[f.Name()] = true
	}

	return fileList, nil
}
