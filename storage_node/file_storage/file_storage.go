package file_storage

import (
	up "Awesome-DFS/transfer"
	"crypto/sha256"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"os"
	"time"
)

type uploadServer struct {
	up.UnimplementedFileTransferServer
}

func (s *uploadServer) Upload(stream up.FileTransfer_UploadServer) error {
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return fmt.Errorf("failed to get peer from context")
	}
	log.Printf("Received upload request from %s\n", p.Addr)

	hasher := sha256.New()
	start := time.Now()
	var chunkFile *os.File
	for {
		chunk, err := stream.Recv()

		if err == io.EOF {
			elapsed := time.Since(start)
			checksum := fmt.Sprintf("%x", hasher.Sum(nil))
			log.Printf("Upload completed in %v with checksum %s\n", elapsed, checksum)
			return stream.SendAndClose(&up.UploadResponse{
				Status:  up.Status_STATUS_OK,
				Message: fmt.Sprintf("Upload completed in %v", elapsed),
			})
		}

		if err != nil {
			return fmt.Errorf("failed to receive chunk metadata: %v", err)
		}

		switch payload := chunk.Payload.(type) {
		case *up.Chunk_Meta:
			chunkFile, err = initChunk(payload.Meta)
			if err != nil {
				return fmt.Errorf("failed to initialize chunk: %v", err)
			}
		case *up.Chunk_Data:
			_, err = chunkFile.Write(payload.Data.RawBytes)
			if err != nil {
				log.Printf("failed to write chunk: %v", err)
				return fmt.Errorf("failed to write chunk: %v", err)
			}
			_, err = hasher.Write(payload.Data.RawBytes)
			if err != nil {
				log.Printf("failed to write chunk to hasher: %v", err)
				return fmt.Errorf("failed to write chunk to hasher: %v", err)
			}
		}
	}
}

func initChunk(meta *up.MetaData) (*os.File, error) {
	filePath := fmt.Sprintf("storage/%s.chunk", meta.UniqueName)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err

	}
	return file, nil
}

func RegisterFileTransferServer(server *grpc.Server) {
	up.RegisterFileTransferServer(server, new(uploadServer))
}
