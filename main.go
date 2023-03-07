package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func searchFiles(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
func main() {
	// Read directory path and file extension from command-line arguments
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s dir ext\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]
	ext := os.Args[2]

	// Search for files with the given extension in the directory
	files, err := searchFiles(dir, ext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for files: %s\n", err.Error())
		os.Exit(1)
	}

	// Retrieve the files to a new directory
	newDir := "retrieved_files"
	err = os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err.Error())
		os.Exit(1)
	}
	for _, file := range files {
		err = copyFile(file, filepath.Join(newDir, filepath.Base(file)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving file: %s\n", err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("%d files retrieved to %s\n", len(files), newDir)
}

func copyFile(src string, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file on the local file system
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Write the contents of the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
