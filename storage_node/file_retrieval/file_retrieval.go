package file_retrieval

import (
	down "Awesome-DFS/protobuf/download"
	ms "Awesome-DFS/storage_node/metadata_service"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
)

var payloadSize int64 = 2 * 1024 * 1024

type DownloadServer struct {
	down.UnimplementedDownloadServer
}

func (s *DownloadServer) Download(chunkDesc *down.ChunkDesc, stream down.Download_DownloadServer) error {
	if !ms.ChunkExists(chunkDesc.FileName, chunkDesc.ChunkName) {
		return fmt.Errorf("chunk does not exist")
	}

	chunkPath := fmt.Sprintf("storage/%s/%s.chunk", chunkDesc.FileName, chunkDesc.ChunkName)

	chunkFile, err := os.Open(chunkPath)
	if err != nil {
		log.Printf("could not open chunk file: %v", err)
		return fmt.Errorf("could not open chunk file: %v", err)
	}
	defer chunkFile.Close()

	chunkInfo, err := chunkFile.Stat()
	if err != nil {
		log.Printf("could not get chunk file info: %v", err)
		return fmt.Errorf("could not get chunk file info: %v", err)
	}

	if chunkInfo.Size() != chunkDesc.Size {
		log.Printf("chunk size mismatch: %d != %d", chunkInfo.Size(), chunkDesc.Size)
		return fmt.Errorf("chunk size mismatch: %d != %d", chunkInfo.Size(), chunkDesc.Size)
	}

	limit := chunkInfo.Size()
	rawBytes := make([]byte, payloadSize)
	for offset := int64(0); offset < chunkInfo.Size(); offset += payloadSize {
		if offset+payloadSize > limit {
			rawBytes = rawBytes[:limit-offset]
		}

		_, err := chunkFile.ReadAt(rawBytes, offset)
		if err != nil {
			log.Printf("could not read chunk file: %v", err)
			return fmt.Errorf("could not read chunk file: %v", err)
		}

		data := &down.Data{RawBytes: rawBytes, Number: offset / payloadSize}
		chunkData := &down.Chunk_Data{Data: data}
		payload := &down.Chunk{Payload: chunkData}

		err = stream.Send(payload)
		if err != nil {
			log.Printf("could not send chunk data: %v", err)
			return fmt.Errorf("could not send chunk data: %v", err)
		}
	}

	checksum, _ := ms.GetChecksum(chunkDesc.FileName, chunkDesc.ChunkName)

	check := &down.IntegrityCheck{Checksum: checksum}
	chunkCheck := &down.Chunk_IntegrityCheck{IntegrityCheck: check}
	payload := &down.Chunk{Payload: chunkCheck}

	err = stream.Send(payload)
	if err != nil {
		log.Printf("could not send integrity check: %v", err)
		return fmt.Errorf("could not send integrity check: %v", err)
	}

	return nil
}

func RegisterDownloadServer(server *grpc.Server) {
	down.RegisterDownloadServer(server, &DownloadServer{})
}
