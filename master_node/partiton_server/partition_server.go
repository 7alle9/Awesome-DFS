package partiton_server

import (
	"Awesome-DFS/master_node/comms_master"
	ms "Awesome-DFS/master_node/metadata_service"
	pb "Awesome-DFS/partition"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"time"
)

type partitionServer struct {
	pb.UnimplementedPartitionServer
}

func (*partitionServer) Split(ctx context.Context, file *pb.File) (*pb.FilePartition, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}

	log.Printf("Partition request for new file from %s\n", p.Addr.String())

	if ms.FileExists(file.Name) {
		log.Printf("Error: File %s already exists", file.Name)
		return nil, fmt.Errorf("file  %s already exists", file.Name)
	}

	addressBook := ms.GetAddressBook()

	availableNodesResponse := comms_master.GetAvailableNodes(addressBook, int(file.ChunkSize))
	if len(availableNodesResponse) == 0 {
		return nil, fmt.Errorf("no available nodes")
	}

	availableNodes := extractAvailableNodes(availableNodesResponse)

	fileUuid := uuid.New().String()

	chunks := chunksInit(fileUuid, file.Size, file.ChunkSize)

	createReplicaChains(chunks, availableNodes, file.NbReplicas)

	choseChainHeads(chunks, availableNodesResponse)

	filePartition := &pb.FilePartition{Chunks: chunks}

	err := ms.StoreFile(fileUuid, file.Name, file.Size, file.ChunkSize, int(file.NbReplicas), filePartition)
	if err != nil {
		return nil, err
	}

	return filePartition, nil
}

func extractAvailableNodes(nodesResponse map[string]time.Duration) []string {
	var nodes []string
	for node, _ := range nodesResponse {
		nodes = append(nodes, node)
	}
	return nodes
}

func chunksInit(fileUuid string, fileSize int64, chunkSize int64) []*pb.Chunk {
	nbChunks := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		nbChunks++
	}

	chunks := make([]*pb.Chunk, nbChunks)
	for i := int64(0); i < nbChunks; i++ {
		name := fmt.Sprintf("%s_chunk_%d", fileUuid, i)
		size := min(fileSize-i*chunkSize, chunkSize)
		offset := i * chunkSize

		chunk := &pb.Chunk{Name: name, Size: size, Offset: offset}
		chunks[i] = chunk
	}

	return chunks
}

func createReplicaChains(chunks []*pb.Chunk, availableNodes []string, nbReplicas int32) {
	curNode := 0

	for _, chunk := range chunks {
		for i := int32(0); i < nbReplicas; i++ {
			chunk.ReplicaChain = append(chunk.ReplicaChain, availableNodes[curNode])
			curNode = (curNode + 1) % len(availableNodes)
		}
	}
}

func choseChainHeads(chunks []*pb.Chunk, nodesResponse map[string]time.Duration) {
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
	pb.RegisterPartitionServer(server, new(partitionServer))
}
