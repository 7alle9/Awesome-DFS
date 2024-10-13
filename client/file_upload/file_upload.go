package file_upload

import (
	fp "Awesome-DFS/client/file_partition"
	part "Awesome-DFS/partition"
	"os"
)

var payloadSize int64 = 2 * 1024 * 1024

func readData(file *os.File, offset int64, data []byte) {
	_, err := file.ReadAt(data, offset)
	if err != nil {
		panic(err)
	}
}

func worker(jobs <-chan *part.Chunk, file *os.File) {
	for chunk := range jobs {
		data := make([]byte, payloadSize)
		for i := chunk.Offset; i < chunk.Offset+chunk.Size; i += payloadSize {
			if i+payloadSize > chunk.Offset+chunk.Size {
				data = make([]byte, chunk.Offset+chunk.Size-i)
			}
			readData(file, i, data)
		}
	}
}

func initJobs(partition *part.FilePartition) chan *part.Chunk {
	jobs := make(chan *part.Chunk, len(partition.Chunks))
	for _, chunk := range partition.Chunks {
		jobs <- chunk
	}
	close(jobs)

	return jobs
}

func UploadFile(file *os.File, chunkSize int64, nbReplicas int) error {
	filePartition, err := fp.GetFilePartition(file, chunkSize, nbReplicas)
	if err != nil {
		return err
	}

	jobs := initJobs(filePartition)

	return nil
}
