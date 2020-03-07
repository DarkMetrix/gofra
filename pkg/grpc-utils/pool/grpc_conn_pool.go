package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"google.golang.org/grpc"
)

//Global connection pool instance
var globalConnectionPool *ConnectionPool = nil

//Get connection pool instance
func GetConnectionPool() *ConnectionPool {
	if globalConnectionPool == nil {
		globalConnectionPool = NewConnectionPool()
	}

	return globalConnectionPool
}

//Connection pool
type ConnectionPool struct {
	mtx sync.Mutex									//Mutex to protect from race condition

	pools map[string]*ConnectionInfo				//Pool map to save all connections
	clientOpts []grpc.DialOption

	maxConnectionPerAddr int						//Max connections for each address
}

//Connection info
type ConnectionInfo struct {
	Conns []*grpc.ClientConn						//Connections
	Index int64										//Index of the next connection
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool {
		pools: make(map[string]*ConnectionInfo),
		clientOpts: make([]grpc.DialOption, 0),
		maxConnectionPerAddr: runtime.NumCPU(),
	}
}

//Init connection pool
func (connPool *ConnectionPool) Init(clientOpts []grpc.DialOption) error {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connPool.clientOpts = clientOpts

	return nil
}

//Close connection pool
func (connPool *ConnectionPool) Close() {
	for _, pool := range connPool.pools {
		if pool == nil {
			continue
		}

		for _, connection := range pool.Conns {
			if connection == nil {
				continue
			}

			connection.Close()
			connection = nil
		}

		pool = nil
	}
}

//Get connection from pool
func (connPool *ConnectionPool) GetConnection(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connInfo, ok := connPool.pools[addr]

	if !ok {
		//Init connection info
		connInfo = &ConnectionInfo{
			Conns: make([]*grpc.ClientConn, connPool.maxConnectionPerAddr),
			Index: 0,
		}

		connPool.pools[addr] = connInfo
	}

	curIndex := connInfo.Index % int64(len(connInfo.Conns))

	//Get connection
	if connInfo.Conns[curIndex] != nil {
		connInfo.Index++
		return connInfo.Conns[curIndex], nil
	} else {
		//Get grpc Client connection
		conn, err := getClientConn(addr, connPool.clientOpts)

		if err != nil {
			return nil, err
		}

		connInfo.Conns[curIndex] = conn
		connInfo.Index++
		return conn, err
	}
}

func getClientConn(addr string, clientOpts []grpc.DialOption) (*grpc.ClientConn, error) {
	// dial remote server
	clientOpts = append(clientOpts, grpc.WithInsecure())

	conn, err := grpc.Dial(addr, clientOpts ...)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("GRPC Dial failed! error:%v", err.Error()))
	}

	return conn, nil
}