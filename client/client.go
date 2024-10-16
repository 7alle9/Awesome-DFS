package main

import (
	fdown "Awesome-DFS/client/file_download"
	fup "Awesome-DFS/client/file_upload"
	"log"
	"os"
)

func main() {
	f, err := os.Open("bigF.txt")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = fup.UploadFile(f, 64*1024*1024, 3)
}

func download() {
	fdown.Download("bigF.txt")
}
