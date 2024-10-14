package metadata_service

import "sync"

var (
	checksums = make(map[string]map[string]string)
	mu        sync.Mutex
)

func newFile(uuid string) {
	checksums[uuid] = make(map[string]string)
}

func NewChunk(fileUuid string, name string, checksum string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := checksums[fileUuid]; !exists {
		newFile(fileUuid)
	}
	checksums[fileUuid][name] = checksum
}

func GetChecksum(fileUuid string, name string) string {
	return checksums[fileUuid][name]
}
