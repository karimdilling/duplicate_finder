package main

import (
	"crypto/sha256"
	"duplicate_finder/options"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 1. Walk filetree and hash all files (excluding folders)
// 2. Save these hashes and their filepath to a map: map[hash]paths
// 3. Print all entries of the map with more than one path
func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("A folder to search in must be entered. Use '.' for the current directory.")
		return
	}

	subDirsToSkip := strings.Split(*options.ExcludeFlag, " ")

	for i := range flag.Args() {
		hash_paths := make(map[string][]string)
		abs_path, err := filepath.Abs(flag.Arg(i))
		if err != nil {
			abs_path = flag.Arg(i)
		}
		err = filepath.WalkDir(abs_path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Error traversing the path: %v\n", err)
				return nil
			}

			skipDir := false
			if d.IsDir() {
				for _, dir := range subDirsToSkip {
					if d.Name() == dir {
						skipDir = true
						break
					}
					skipDir = false
				}
			}
			if d.IsDir() && skipDir {
				return filepath.SkipDir
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
