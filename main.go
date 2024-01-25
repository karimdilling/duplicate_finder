package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// 1. Walk filetree and hash all files (excluding folders)
// 2. Save these hashes and their filepath to a map: map[hash]paths
// 3. Print all entries of the map with more than one path
func main() {
	if len(os.Args) < 2 {
		fmt.Println("A folder to search in must be entered. Use '.' for the current directory.")
		return
	}

	for i := range os.Args {
		if i == 0 {
			continue
		}
		hash_paths := make(map[string][]string)
		fmt.Printf("############################### Duplicates in path %v ######################################\n", os.Args[i])
		err := filepath.WalkDir(os.Args[i], func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Error traversing the path: %v\n", err)
				return nil
			}

			if d.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("Could not open file: %v\n", err)
				return nil
			}
			defer file.Close()
			hash, err := calcHashForFile(file)
			if err != nil {
				fmt.Printf("Could not generate unique identifier for file: %v\n", err)
				return nil
			}
			paths := append(hash_paths[hash], path)
			hash_paths[hash] = paths
			return nil
		})
		if err != nil {
			fmt.Printf("Access error: %v\n", err)
		}

		printDuplicates(hash_paths)
	}

}

func calcHashForFile(file *os.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func printDuplicates(hash_paths map[string][]string) {
	for _, paths := range hash_paths {
		if len(paths) > 1 {
			for i, path := range paths {
				if i < len(paths)-1 {
					fmt.Printf("%d: %v\n", i+1, path)
				} else {
					fmt.Printf("%d: %v\n\n", i+1, path)
				}
			}
		}
	}
}
