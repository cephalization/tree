package main

import (
	"fmt"
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
		err := tree(arg)
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}
}

func tree(root string) error {
	// Walk each file starting at root
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()

		// don't print "." for current dir
		if strings.Compare(name, ".") == 0 {
			return nil
		}

		// skip walking hidden directories but do not skip relative pathing
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "..") {
			return filepath.SkipDir
		}

		fmt.Println(name)

		return nil
	})

	return err
}
