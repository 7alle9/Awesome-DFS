package chunk_forwarding

import (
	"Awesome-DFS/protobuf/transfer"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"os"
)

var (
	opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10, Timeout: 5, PermitWithoutStream: true}),
	}
	payloadSize int64 = 2 * 1024 * 1024
)

func nextNode(chain *[]string) string {
	if len(*chain) == 0 {
		return ""
	}
	res := (*chain)[0]
	*chain = (*chain)[1:]
	return res
}

func getStream(addr string) (
	grpc.ClientStreamingClient[__.Chunk, __.UploadResponse],
	*grpc.ClientConn,
	error,
) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, nil, err
	}

	client := __.NewFileTransferClient(conn)
	stream, err := client.Upload(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return stream, conn, nil
}

func readData(file *os.File, offset int64, data []byte) {
	_, err := file.ReadAt(data, offset)
	if err != nil {
		panic(err)
	}
}

func Next(chunkFile *os.File, metadata *__.MetaData) {
	defer chunkFile.Close()

	forwardTo := nextNode(&metadata.ReplicaChain)
	if forwardTo == "" {
		log.Printf("Chunk %s reached end of chain", metadata.UniqueName)
		return
	}
	log.Printf("Forwarding chunk %s to %s", metadata.UniqueName, forwardTo)

	stream, conn, err := getStream(forwardTo)
	if err != nil {
		log.Printf("failed to get stream to forward chunk: %v", err)
		return
	}
	defer conn.Close()

	chunkMeta := &__.Chunk_Meta{Meta: metadata}
	chunk := &__.Chunk{Payload: chunkMeta}

	err = stream.Send(chunk)
	if err != nil {
		log.Printf("error sending metadata: %v", err)
		return
	}

	data := make([]byte, payloadSize)
	limit := metadata.Size
	for i := int64(0); i < limit; i += payloadSize {
		if i+payloadSize > limit {
			data = data[:limit-i]
		}

		readData(chunkFile, i, data)

		payloadData := &__.Data{RawBytes: data, Number: i / payloadSize}
		chunkData := &__.Chunk_Data{Data: payloadData}
		chunk.Payload = chunkData

		err = stream.Send(chunk)
		if err != nil {
			log.Printf("error sending data: %v", err)
			return
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("failed to close and receive: %v", err)
		return
	}

	if reply.Status == __.Status_STATUS_OK {
		log.Printf("Chunk %s forwarded successfully", metadata.UniqueName)
	} else {
		log.Printf("failed to forward chunk %s: %s", metadata.UniqueName, reply.Message)
	}
}
