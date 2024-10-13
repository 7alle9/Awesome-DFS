package file_upload

import (
	fp "Awesome-DFS/client/file_partition"
	part "Awesome-DFS/partition"
	"os"
)

func UploadFile(file *os.File, chunkSize int64, nbReplicas int) error {
	filePartition, err := fp.GetFilePartition(file, chunkSize, nbReplicas)
	if err != nil {
		return err
	}

	jobs := initJobs(filePartition)

	return nil
}

func initJobs(partition *part.FilePartition) chan *part.Chunk {
	jobs := make(chan *part.Chunk, len(partition.Chunks))
	for _, chunk := range partition.Chunks {
		jobs <- chunk
	}
	close(jobs)

	return jobs
}
