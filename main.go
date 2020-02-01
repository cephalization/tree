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

	// start building tree for each root dir
	for _, arg := range args {
		err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}
}

func tree(root, indent string) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("Could not stat %s: %v", root, err)
	}

	// Print initial indent + additional indent then the file name
	fmt.Printf("%s%s\n", indent, info.Name())

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
	for _, info := range infos {
		// Skip hidden files and directories
		// Additionally skips 'node_modules' but this will be replaced with blacklist map soon
		if info.Name()[0] == '.' || strings.Compare(info.Name(), "node_modules") == 0 {
			continue
		}

		cwd := filepath.Join(root, info.Name())
		newIndent := indent + "  "
		err := tree(cwd, newIndent)
		if err != nil {
			return fmt.Errorf("Could not print tree at %s: %v", cwd, err)
		}
	}

	return nil
}
