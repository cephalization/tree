package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Default tree to current dir
	args := []string{"."}

	// Use provided root directories if available
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	files := make([]os.FileInfo, len(args))

	for i, arg := range args {
		file, err := os.Stat(arg)
		if err != nil {
			log.Fatalf("Could not stat %s: %v", arg, err)
			return
		}

		files[i] = file
	}

	// start building tree for each root dir
	for i, arg := range args {
		err := tree(arg, "", i, files)
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}
}

func tree(root, indent string, index int, lastDir []os.FileInfo) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("Could not stat %s: %v", root, err)
	}

	// Print initial indent + additional indent then the file name
	pipe := "├──"
	if lastDir != nil && index == len(lastDir)-1 {
		pipe = "└──"
	}

	fmt.Printf("%s%s%s\n", indent, pipe, info.Name())

	// Bail out when the file is not a directory
	if !info.IsDir() {
		return nil
	}

	// Read the files in this directory
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("Could not read dir %s: %v", root, err)
	}

	// Recursively call tree with new indentation on each file in directory
	for i, newInfo := range infos {
		// Skip hidden files and directories
		// Additionally skips 'node_modules' but this will be replaced with blacklist map soon
		if newInfo.Name()[0] == '.' || strings.Compare(newInfo.Name(), "node_modules") == 0 {
			continue
		}

		// Print tree characters
		//│
		//├─
		//└─
		cwd := filepath.Join(root, newInfo.Name())
		newIndent := indent + "│  "
		if info.IsDir() && index == len(lastDir)-1 {
			newIndent = indent + "   "
		}

		err := tree(cwd, newIndent, i, infos)
		if err != nil {
			return fmt.Errorf("Could not print tree at %s: %v", cwd, err)
		}
	}

	return nil
}
