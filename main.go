package main

import (
	m "Awesome-DFS/model"
	"flag"
	"fmt"
	"log"
)

var fileFlag = flag.String("file", "", "File to fetch")

func main() {
	//flag.Parse()
	f, err := m.FetchLocalFile("file.txt")
	if err != nil {
		log.Fatalln("Error fetching file: ", err)
	}

	if err := f.SplitFile(2); err != nil {
		log.Fatalln("Error splitting file: ", err)
	}
	fmt.Println(f)

	for _, chunk := range f.Chunks {
		fmt.Println(chunk)
	}

	fClone := &m.File{
		Name:     "fileclone.txt",
		Size:     f.Size,
		CheckSum: f.CheckSum,
		Chunks:   f.Chunks,
	}
	if err := fClone.ReconstructFile(); err != nil {
		log.Fatalln("Error reconstructing file: ", err)
	}
	fmt.Println(fClone)
	if err := fClone.WriteFile(); err != nil {
		log.Fatalln("Error writing file: ", err)
	}
}
