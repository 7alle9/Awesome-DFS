package metadata_service

import (
	pb "Awesome-DFS/protobuf/partition"
	"fmt"
	"log"
	"sync"
	"time"
)

type File struct {
	Uuid              string
	FileName          string
	Size              int64
	ChunkSize         int64
	ReplicationFactor int
	LastUpdated       time.Time
	Partition         *pb.FilePartition
}

type TempFile struct {
	File
	ValidationFactor int
	ValidationTarget int
}

type StoreNode struct {
	Ip   string
	Port int
}

var (
	mu         = &sync.Mutex{}
	files      = make(map[string]*File)
	tempFiles  = make(map[string]*TempFile)
	storeNodes = make([]*StoreNode, 0)
)

func FileExists(fileName string) bool {
	_, exists := files[fileName]
	return exists
}

func StoreFile(
	uuid string,
	fileName string,
	size int64,
	chunkSize int64,
	replicationFactor int,
	partition *pb.FilePartition,
) error {
	if FileExists(fileName) {
		return fmt.Errorf("file  %s already exists with uuid %s", fileName, files[fileName].Uuid)
	}

	newFile := &File{
		uuid,
		fileName,
		size,
		chunkSize,
		replicationFactor,
		time.Now(),
		partition,
	}

	files[fileName] = newFile

	return nil
}

func UploadRequest(
	uuid string,
	fileName string,
	size int64,
	chunkSize int64,
	replicationFactor int,
	partition *pb.FilePartition,
) error {
	newTempFile := &TempFile{
		File: File{
			uuid,
			fileName,
			size,
			chunkSize,
			replicationFactor,
			time.Now(),
			partition,
		},
		ValidationFactor: 0,
		ValidationTarget: replicationFactor * len(partition.Chunks),
	}

	tempFiles[uuid] = newTempFile

	return nil
}

func Validate(fileUuid string) {
	mu.Lock()
	defer mu.Unlock()

	tempFile, exists := tempFiles[fileUuid]
	if !exists {
		return
	}

	tempFile.ValidationFactor++
	if tempFile.ValidationFactor == tempFile.ValidationTarget {
		delete(tempFiles, fileUuid)
		files[tempFile.FileName] = &tempFile.File

		log.Printf("File %s validated\n", tempFile.FileName)
	}
}

func GetFile(fileName string) (*File, error) {
	mu.Lock()
	defer mu.Unlock()

	file, exists := files[fileName]
	if !exists {
		return nil, fmt.Errorf("file %s does not exist", fileName)
	}

	return file, nil
}

func (n *StoreNode) Addr() string {
	return fmt.Sprintf("%s:%d", n.Ip, n.Port)
}

func GetAddressBook() []string {
	addresses := make([]string, 0)
	for _, node := range storeNodes {
		addresses = append(addresses, node.Addr())
	}
	return addresses
}

func init() {
	storeNodes = append(storeNodes, &StoreNode{"localhost", 8080})
	storeNodes = append(storeNodes, &StoreNode{"localhost", 8081})
	storeNodes = append(storeNodes, &StoreNode{"localhost", 8082})
	storeNodes = append(storeNodes, &StoreNode{"localhost", 8083})
}
