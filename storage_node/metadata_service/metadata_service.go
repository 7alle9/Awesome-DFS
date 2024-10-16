package metadata_service

import (
	"fmt"
	"sync"
)

var (
	checksums = make(map[string]map[string]string)
	mu        sync.Mutex
)

func newFile(uuid string) {
	checksums[uuid] = make(map[string]string)
}

func ChunkExists(fileUuid string, uniqueName string) bool {
	_, exists := checksums[fileUuid][uniqueName]
	return exists
}

func NewChunk(fileUuid string, name string, checksum string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := checksums[fileUuid]; !exists {
		newFile(fileUuid)
	}
	checksums[fileUuid][name] = checksum
}

func GetChecksum(fileUuid string, uniqueName string) (string, error) {
	if ChunkExists(fileUuid, uniqueName) {
		return checksums[fileUuid][uniqueName], nil
	}
	return "", fmt.Errorf("chunk does not exist")
}
