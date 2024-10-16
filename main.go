package main

import "os"

func main() {
	// create 2gb file
	data := make([]byte, 2*1024*1024*1024)
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < 32; i++ {
		for j := 0; j < 67108864; j++ {
			data[i*67108864+j] = alphabet[i%26]
		}
	}

	file, _ := os.Create("bigF.txt")
	file.Write(data)
}
