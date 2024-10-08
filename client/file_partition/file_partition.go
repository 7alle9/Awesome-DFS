package file_partition

import (
	servers "Awesome-DFS/client/server_connection"
	pb "Awesome-DFS/partition"
	"context"
	"log"
	"os"
)

func getFileSize(file *os.File) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func GetFilePartition(file *os.File, chunkSize int64, nbReplicas int) (*pb.FilePartition, error) {
	partitionServer := servers.GetPartitionClient()

	fileName := file.Name()
	fileSize, err := getFileSize(file)
	if err != nil {
		return nil, err
	}

	splitDescription := &pb.File{
		Name:       fileName,
		Size:       fileSize,
		ChunkSize:  chunkSize,
		NbReplicas: int32(nbReplicas),
	}

	log.Printf("Requesting partition for file %s\n", fileName)

	partition, err := partitionServer.Split(context.Background(), splitDescription)

	return partition, nil
}
