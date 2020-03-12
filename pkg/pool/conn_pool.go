package pool

import (
	"context"
	"net"
	"sync"
	"time"

	commonPool "github.com/silenceper/pool"
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

	pools map[string]commonPool.Pool				// pool map to save all connections

	initConnections int								// initial connection count per addr
	maxConnections int								// max connection count per addr
	idleTimeout time.Duration						// idle timeout for connection
}

// conn is the wrapper for a net.Conn
type Conn struct {
	net.Conn
	connPool  commonPool.Pool
	unhealthy bool
}

// get returns the real connection to use
func (conn *Conn) Get() net.Conn {
	return conn.Conn
}

// unhealthy mark the connection as unhealthy
// when recycle called it will be closed and won't be put back to the pool
func (conn *Conn) Unhealthy() {
	conn.unhealthy = true
}

// recycle returns the connection to the pool
// if the unhealthy mark is set, close and it won't be put back to the pool
func (conn *Conn) Recycle() error {
	if conn.unhealthy {
		conn.connPool.Close(conn)
		return nil
	} else {
		err := conn.connPool.Put(conn)

		if err != nil {
			return err
		}

		return nil
	}
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool {
		pools: make(map[string]commonPool.Pool),
		initConnections: 10,
		maxConnections: 500,
		idleTimeout: time.Second * 30,
	}
}

// init connection pool
func (connPool *ConnectionPool) Init(initConnections, maxConnections int, idleTimeout time.Duration) error {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	connPool.initConnections = initConnections
	connPool.maxConnections = maxConnections
	connPool.idleTimeout = idleTimeout

	return nil
}

// get connection from pool
func (connPool *ConnectionPool) GetConnection(ctx context.Context, addr string) (*Conn, error) {
	connPool.mtx.Lock()
	defer connPool.mtx.Unlock()

	var err error
	pool, ok := connPool.pools[addr]

	if !ok {
		pool, err = commonPool.NewChannelPool(&commonPool.Config{
			InitialCap: connPool.initConnections,
			MaxCap: connPool.maxConnections,
			Factory: getFactory(addr),
			Close: closeFunc,
			IdleTimeout: connPool.idleTimeout,
		})

		if err != nil {
			return nil, err
		}

		connPool.pools[addr] = pool
	}

	wrapper := &Conn {
		connPool: pool,
		unhealthy: false,
	}

	conn, err := pool.Get()

	if err != nil {
		return nil, err
	}

	wrapper.Conn = conn.(net.Conn)

	return wrapper, nil
}

func getFactory(addr string) func() (interface{}, error) {
	return func() (interface{}, error) {
		return net.Dial("tcp", addr)
	}
}

func closeFunc(conn interface{}) error {
	var err error = nil

	if conn != nil {
		err = conn.(net.Conn).Close()
	}

	return err
}
