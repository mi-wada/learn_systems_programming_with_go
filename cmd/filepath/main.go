package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	path := filepath.Join(os.TempDir(), "test.dir")
	fmt.Printf("path: %v\n", path)
	dir, file := filepath.Split(path)
	fmt.Printf("dir: %v, file: %v\n", dir, file)
	parts := filepath.SplitList(os.Getenv("PATH"))
	fmt.Printf("parts: %v\n", parts)

	dirtyPath := "../hoge/../hoge/fuga"
	cleanPath := filepath.Clean(dirtyPath)
	fmt.Printf("dirtyPath: %v, cleanPath: %v\n", dirtyPath, cleanPath)

	fmt.Println(os.UserHomeDir())

	files, err := filepath.Glob("test.*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("files: %v\n", files)
}
