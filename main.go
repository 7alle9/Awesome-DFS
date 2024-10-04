package main

import (
	"fmt"
	"os"
)

func main() {
	// current path
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current path: %s\n", dir)
}
