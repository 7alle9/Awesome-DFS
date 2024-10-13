package main

import "os"

func main() {
	testData := make([]byte, 2*1024*1024)
	f, err := os.Create("bigF.txt")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 2*1024*1024; i++ {
		testData[i] = 'a'
	}

	_, err = f.Write(testData)
	if err != nil {
		panic(err)
	}
}
