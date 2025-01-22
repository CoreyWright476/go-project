package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var IgnoreDir = []string{".git", "loki-build-scripts"}

func main() {
	path := getPathFromArgs()

	fullPath := resolveFullPath(path)
	log.Printf("Got path: %s\n", fullPath)

	crawlPath(fullPath)

	fmt.Printf("RESULT - Directory Count: %d\t\tFile Count: %d\n", dirCount, fileCount)
	fmt.Println(directoryToFileCount)
}

// path must be first arg
// Assuming path provided is from the home dir
func getPathFromArgs() string {
	if len(os.Args) < 2 {
		log.Fatalf("No path provided")
	}

	return os.Args[1]
}

// prepend home to user provided path
func resolveFullPath(userProvidedPath string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Couldn't resolve user home directory")
	}
	return homeDir + userProvidedPath
}

var callCount = 0

var dirCount = 0
var fileCount = 0
var directoryToFileCount = make(map[string]int)

func crawlPath(directory string) {
	filepath.WalkDir(directory, visit)
}

func visit(path string, d os.DirEntry, error error) error {
	callCount++

	//If directory is an ignored directory - skip it
	if d.IsDir() {
		if slices.Contains(IgnoreDir, d.Name()) {
			return filepath.SkipDir
		}
		dirCount++
		// Add to map to count per dir?
		if _, ok := directoryToFileCount[d.Name()]; ok {
			directoryToFileCount[d.Name()] = 0
		}
	} else {
		fileCount++
		//get dir of this file,
		var strippedPath = strings.TrimSuffix(path, d.Name())
		//add 1 to value of map where key is this dir
		directoryToFileCount[strippedPath]++
	}

	fmt.Printf("Call no: %d\tPath: %s\t\n", callCount, path)
	return nil
}
