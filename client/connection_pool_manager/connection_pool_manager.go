package connection_pool_manager

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"sync"
)

var (
	mu            sync.Mutex
	idIncrement   = 0
	connection    = make(map[int]*grpc.ClientConn)
	workerPerConn = make(map[int]int)
	connAddr      = make(map[int]string)
	connPool      = make(map[string][]int)
	opts          = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10, Timeout: 5, PermitWithoutStream: true}),
	}
)

func newConn(address string) (*grpc.ClientConn, int) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		panic(err)
	}

	mu.Lock()

	connId := idIncrement
	idIncrement++
	workerPerConn[connId] = 1
	connAddr[connId] = address
	connection[connId] = conn
	connPool[address] = append(connPool[address], connId)

	mu.Unlock()

	log.Printf("New connection to %s created with id %d", address, connId)

	return conn, connId
}

func isSaturated(connId int) bool {
	return workerPerConn[connId] >= 50
}

func ConnectTo(address string) (*grpc.ClientConn, int) {
	mu.Lock()

	if len(connPool[address]) == 0 {
		mu.Unlock()

		return newConn(address)
	} else {
		mu.Unlock()
	}

	mu.Lock()

	connID := connPool[address][0]
	workerPerConn[connID]++
	if isSaturated(connID) {
		connPool[address] = connPool[address][1:]
		log.Printf("Connection to %s with id %d is saturated", address, connID)
	}

	mu.Unlock()

	return connection[connID], connID
}

func ReleaseConn(connID int) {
	mu.Lock()

	if isSaturated(connID) {
		connPool[connAddr[connID]] = append(connPool[connAddr[connID]], connID)
		log.Printf("Connection to %s with id %d is no longer saturated", connAddr[connID], connID)
	}
	workerPerConn[connID]--

	mu.Unlock()
}

func CloseAll() {
	clear(connAddr)
	clear(connPool)
	clear(connection)
	clear(workerPerConn)
	idIncrement = 0
}
