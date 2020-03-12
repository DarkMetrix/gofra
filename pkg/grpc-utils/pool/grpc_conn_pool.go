package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"google.golang.org/grpc"
)

// global connection pool instance
var globalConnectionPool *ConnectionPool = nil

// get connection pool instance
func GetConnectionPool() *ConnectionPool {
	if globalConnectionPool == nil {
		globalConnectionPool = NewConnectionPool()
	}

	return globalConnectionPool
}

// connection pool
type ConnectionPool struct {
	mtx sync.Mutex									// mutex to protect from race condition

	pools map[string]*ConnectionInfo				// pool map to save all connections
	clientOpts []grpc.DialOption

	maxConnectionPerAddr int						// max connections for each address
}

// connection info
type ConnectionInfo struct {
	Conns []*grpc.ClientConn						// connections
	Index int64										// index of the next connection
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool {
		pools: make(map[string]*ConnectionInfo),
		clientOpts: make([]grpc.DialOption, 0),
		maxConnectionPerAddr: runtime.NumCPU(),
	}
}

// init connection pool
func (connPool *ConnectionPool) Init(clientOpts []grpc.DialOption) error {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connPool.clientOpts = clientOpts

	return nil
}

// close connection pool
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

// get connection from pool
func (connPool *ConnectionPool) GetConnection(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connInfo, ok := connPool.pools[addr]

	if !ok {
		// init connection info
		connInfo = &ConnectionInfo{
			Conns: make([]*grpc.ClientConn, connPool.maxConnectionPerAddr),
			Index: 0,
		}

		connPool.pools[addr] = connInfo
	}

	curIndex := connInfo.Index % int64(len(connInfo.Conns))

	// get connection
	if connInfo.Conns[curIndex] != nil {
		connInfo.Index++
		return connInfo.Conns[curIndex], nil
	} else {
		// get grpc Client connection
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