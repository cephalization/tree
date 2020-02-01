package main

import (
	"fmt"
	"os"
)

func main() {
	args := []string{"."}

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	fmt.Println(args)
}
