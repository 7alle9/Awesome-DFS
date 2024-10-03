package main

import (
	pb "Awesome-DFS/dispatcher"
	"context"
	"fmt"
	"log"
	"os"
)

//var fileFlag = flag.String("file", "", "File to fetch")

func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	handleError(err, "close error")
}

func createChunk(data []byte, number int64) {
	chunk, err := os.Create(fmt.Sprintf("chunks/chunk_%d", number))
	handleError(err, "create error")
	defer closeFile(chunk)

	_, err = chunk.Write(data)
	handleError(err, "write error")
}

func readChunk(ind int64) []byte {
	chunk, err := os.Open(fmt.Sprintf("chunks/chunk_%d", ind))
	handleError(err, "open error")
	defer closeFile(chunk)

	info, err := chunk.Stat()
	handleError(err, "path error")

	data := make([]byte, info.Size())
	_, err = chunk.Read(data)
	handleError(err, "read error")

	return data
}

func splitFile(filename string) {
	f, err := os.Open(filename)
	handleError(err, "path error")
	defer closeFile(f)

	info, err := f.Stat()
	handleError(err, "path error")

	size := info.Size()
	log.Printf("Size: %d\n", size)
	var chunkSize int64 = 64 * 1024

	buf := make([]byte, chunkSize)
	for off := int64(0); off < size; off += chunkSize {
		if off+chunkSize > size {
			buf = buf[:size-off]
		}
		_, err := f.ReadAt(buf, off)
		handleError(err, "read error")

		createChunk(buf, off/chunkSize)
	}
	log.Printf("Done splitting: %s\n", filename)
}

func reconstructFile(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err, "file creation error")
	defer closeFile(f)

	for i := int64(0); i < 22541; i++ {
		buf := readChunk(i)

		_, err = f.Write(buf)
		handleError(err, "file write error")
	}
}

type StorageNode struct {
	pb.UnimplementedDispatcherServer
}

func (s *StorageNode) Dispatch(ctx context.Context, in *pb.Chunk) (*pb.Response, error) {
	
}

func main() {

}
