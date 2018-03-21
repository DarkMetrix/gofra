package pool

import (
	"fmt"
	"sync"
	"time"
	"errors"
	"context"

	"google.golang.org/grpc"
	grpcPool "github.com/processout/grpc-go-pool"
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

	pools map[string]*grpcPool.Pool					//Pool map to save all connections
	clientOpts []grpc.DialOption

	initConnections int								//Initial connection count per addr
	maxConnections int								//Max connection count per addr
	idleTimeout time.Duration						//Idle timeout for connection
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool {
		pools: make(map[string]*grpcPool.Pool),
		clientOpts: make([]grpc.DialOption, 0),
		initConnections: 10,
		maxConnections: 500,
		idleTimeout: time.Second * 30,
	}
}

//Init connection pool
func (connPool *ConnectionPool) Init(clientOpts []grpc.DialOption,
	initConnections, maxConnections int, idleTimeout time.Duration) error {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connPool.clientOpts = clientOpts
	connPool.initConnections = initConnections
	connPool.maxConnections = maxConnections
	connPool.idleTimeout = idleTimeout

	return nil
}

//Get connection from pool
func (connPool *ConnectionPool) GetConnection(ctx context.Context, addr string) (*grpcPool.ClientConn, error) {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	var err error
	pool, ok := connPool.pools[addr]

	if !ok {
		pool, err = grpcPool.New(getFactory(addr, connPool.clientOpts), connPool.initConnections, connPool.maxConnections, connPool.idleTimeout)

		if err != nil {
			return nil, err
		}

		connPool.pools[addr] = pool
	}

	conn, err := pool.Get(ctx)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getFactory(addr string, clientOpts []grpc.DialOption) grpcPool.Factory {
	return func() (*grpc.ClientConn, error) {
		// dial remote server
		clientOpts := append(clientOpts, grpc.WithInsecure())

		conn, err := grpc.Dial(addr, clientOpts ...)

		if err != nil {
			fmt.Println(err)
			return nil, errors.New(fmt.Sprintf("GRPC Dial failed! error:%v", err.Error()))
		}

		return conn, nil
	}
}