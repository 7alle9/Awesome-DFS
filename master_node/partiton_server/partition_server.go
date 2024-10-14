package partiton_server

import (
	"Awesome-DFS/master_node/comms_master"
	ms "Awesome-DFS/master_node/metadata_service"
	"Awesome-DFS/protobuf/partition"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"time"
)

type partitionServer struct {
	__.UnimplementedPartitionServer
}

func (*partitionServer) Split(ctx context.Context, file *__.File) (*__.FilePartition, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}

	log.Printf("Upload request for new file from %s\n", p.Addr.String())

	addressBook := ms.GetAddressBook()

	availableNodesResponse := comms_master.GetAvailableNodes(addressBook, int(file.ChunkSize))
	if len(availableNodesResponse) == 0 {
		return nil, fmt.Errorf("no available storage nodes")
	}

	availableNodes := extractAvailableNodes(availableNodesResponse)

	fileUuid := uuid.New().String()

	chunks := chunksInit(fileUuid, file.Size, file.ChunkSize)

	createReplicaChains(chunks, availableNodes, file.NbReplicas)

	choseChainHeads(chunks, availableNodesResponse)

	filePartition := &__.FilePartition{FileUuid: fileUuid, Chunks: chunks}

	err := ms.UploadRequest(fileUuid, file.Name, file.Size, file.ChunkSize, int(file.NbReplicas), filePartition)
	if err != nil {
		return nil, err
	}

	return filePartition, nil
}

func (*partitionServer) Reconstruct(ctx context.Context, desc *__.FileDesc) (*__.FilePartition, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}

	log.Printf("%s request to download %s\n", p.Addr.String(), desc.Filename)

	file, err := ms.GetFile(desc.Filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return file.Partition, nil
}

func extractAvailableNodes(nodesResponse map[string]time.Duration) []string {
	var nodes []string
	for node, _ := range nodesResponse {
		nodes = append(nodes, node)
	}
	return nodes
}

func chunksInit(fileUuid string, fileSize int64, chunkSize int64) []*__.Chunk {
	nbChunks := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		nbChunks++
	}

	chunks := make([]*__.Chunk, nbChunks)
	for i := int64(0); i < nbChunks; i++ {
		name := fmt.Sprintf("%s_chunk_%d", fileUuid, i)
		size := min(fileSize-i*chunkSize, chunkSize)
		offset := i * chunkSize

		chunk := &__.Chunk{Name: name, Size: size, Offset: offset}
		chunks[i] = chunk
	}

	return chunks
}

func createReplicaChains(chunks []*__.Chunk, availableNodes []string, nbReplicas int32) {
	curNode := 0

	for _, chunk := range chunks {
		for i := int32(0); i < nbReplicas; i++ {
			chunk.ReplicaChain = append(chunk.ReplicaChain, availableNodes[curNode])
			curNode = (curNode + 1) % len(availableNodes)
		}
	}
}

func choseChainHeads(chunks []*__.Chunk, nodesResponse map[string]time.Duration) {
	for _, chunk := range chunks {
		chainHeadIndex := 0
		chainHead := chunk.ReplicaChain[chainHeadIndex]

		for i, node := range chunk.ReplicaChain {
			if nodesResponse[node] < nodesResponse[chainHead] {
				chainHead = node
				chainHeadIndex = i
			}
		}

		chunk.SendTo = chainHead
		chunk.ReplicaChain = append(chunk.ReplicaChain[:chainHeadIndex], chunk.ReplicaChain[chainHeadIndex+1:]...)
	}
}

func RegisterPartitionServer(server *grpc.Server) {
	__.RegisterPartitionServer(server, new(partitionServer))
}
