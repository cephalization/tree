package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

func main() {
	// Default tree to current dir
	args := []string{"."}

	// Use provided root directories if available
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	files := make([]os.FileInfo, len(args))

	// validate each argument as a file/directory
	for i, arg := range args {
		file, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		files[i] = file
	}

	// start building tree for each root dir
	for i, arg := range args {
		err := tree(arg, "", i, 0, files)
		if err != nil {
			log.Printf("could not generate tree %s: %v\n", arg, err)
		}
	}
}

func tree(root, indent string, index, depth int, lastDir []os.FileInfo) error {
	file, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("Could not stat %s: %v", root, err)
	}

	// Determine correct tree character
	//│
	//├─
	//└─
	pipe := ""

	pipe = "├──"
	if lastDir != nil && index == len(lastDir)-1 {
		pipe = "└──"
	}

	fileNameFormatted := file.Name()

	// When we are at the top level of the tree command,
	// print the absolute paths of the requested directories to tree
	if depth == 0 {
		if !filepath.IsAbs(root) {
			dir, err := filepath.Abs(fileNameFormatted)
			if err != nil {
				return fmt.Errorf("Could not get absolute path %s:%v", fileNameFormatted, err)
			}

			if root != "." {
				cDir := filepath.Dir(dir)
				fileNameFormatted = filepath.Join(cDir, root)
			} else {
				fileNameFormatted = dir
			}
		} else {
			fileNameFormatted = root
		}

		// Print directory name in green
		color.Green.Println(fileNameFormatted)
	} else {

		// Print file / directory name with decorations
		if file.IsDir() {
			fmt.Printf("%s%s", indent, pipe)
			color.Cyan.Printf("%s\n", fileNameFormatted)
		} else {
			fmt.Printf("%s%s%s\n", indent, pipe, fileNameFormatted)
		}
	}

	// Bail out when the file is not a directory
	if !file.IsDir() {
		return nil
	}

	// Read the files in this directory
	directory, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("Could not read dir %s: %v", root, err)
	}

	// Recursively call tree with new indentation on each file in directory
	for i, nextFile := range directory {
		// Skip hidden files and directories
		// Additionally skips 'node_modules' but this will be replaced with blacklist map soon
		if nextFile.Name()[0] == '.' || strings.Compare(nextFile.Name(), "node_modules") == 0 {
			continue
		}

		// Determine directory indentation
		newIndent := ""
		if depth > 0 {
			newIndent = indent + "│  "
			if file.IsDir() && index == len(lastDir)-1 {
				newIndent = indent + "   "
			}
		}

		// Get the next path
		cwd := filepath.Join(root, nextFile.Name())

		// Recurse
		err := tree(cwd, newIndent, i, depth+1, directory)
		if err != nil {
			return fmt.Errorf("Could not print tree at %s: %v", cwd, err)
		}
	}

	return nil
}
