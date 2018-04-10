package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	log "github.com/cihub/seelog"

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

//ClientConn is the wrapper for a grpc-go-pool.ClientConn
type ClientConn struct {
	*grpcPool.ClientConn
}

//Get returns the real connection to use
func (conn *ClientConn) Get() *grpc.ClientConn {
	return conn.ClientConn.ClientConn
}

//Recycle returns the connection to the pool
//If the unhealthy mark is set, close and it won't be put back to the pool
func (conn *ClientConn) Recycle() error {
	return conn.Close()
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

	log.Tracef("init grpc conn pool success! opts:%v, init conn:%v, max conn:%v, idle timeout:%v", clientOpts, initConnections, maxConnections, idleTimeout)

	return nil
}

//Get connection from pool
func (connPool *ConnectionPool) GetConnection(ctx context.Context, addr string) (*ClientConn, error) {
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

	wrapper := &ClientConn{
		ClientConn:	conn,
	}

	return wrapper, nil
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