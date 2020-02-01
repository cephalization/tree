package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("Could not stat %s: %v", root, err)
	}

	fmt.Println(info.Name())
	if !info.IsDir() {
		return nil
	}

	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("Could not read dir %s: %v", root, err)
	}

	for _, info := range infos {
		cwd := filepath.Join(root, info.Name())
		err := tree(cwd)
		if err != nil {
			return fmt.Errorf("Could not print tree at %s: %v", cwd, err)
		}
	}

	return nil
}
