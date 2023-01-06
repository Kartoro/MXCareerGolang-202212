package connection

import (
	"MXCareerGolang-202212/config"
	"errors"
	"fmt"
	"net"
	"sync"
)

type Pool interface {
	Get() (net.Conn, error)
	Close()
	Remove(conn net.Conn) error
}

// ConnPool implements Pool interface. Use channel buffer connections.
type ConnPool struct {
	lock         sync.Mutex
	connChan     chan net.Conn
	minConnNum   int
	maxConnNum   int
	totalConnNum int
	closed       bool
	connCreator  func() (net.Conn, error)
}

type CpConn struct {
	net.Conn
	pool *ConnPool
}

// NewPool return new Pool. It based on channel. NewPool init minConn connections in channel first.
// When Get()/GetWithTimeout called, if channel still has connection it will get connection from channel.
// Otherwise, ConnPool check number of connection which had already created as the number are less than maxConn,
// it uses connCreator function to create new connection.
func NewPool(minConn, maxConn int, connCreator func() (net.Conn, error)) *ConnPool {
	if minConn > maxConn || minConn < 0 || maxConn <= 0 {
		return nil
	}
	pool := &ConnPool{}
	pool.minConnNum = minConn
	pool.maxConnNum = maxConn
	pool.connCreator = connCreator
	pool.connChan = make(chan net.Conn, maxConn)
	pool.closed = false
	pool.totalConnNum = 0
	err := pool.init()
	if err != nil {
		return nil
	}
	return pool
}

func ConnCreator() (net.Conn, error) {
	return net.Dial("tcp", config.TCPServerAddr)
}

func InitialisePoolValue(pool *ConnPool) *ConnPool {
	if pool == nil {
		fmt.Println("Initialised pool")
		return NewPool(config.TCPClientPoolSize, config.TCPClientPoolSize, ConnCreator)
	}
	return pool
}

func (p *ConnPool) init() error {
	for i := 0; i < p.minConnNum; i++ {
		conn, err := p.createConn()
		if err != nil {
			return err
		}
		p.connChan <- conn
	}
	return nil
}

// Get get connection from connection pool. If connection poll is empty and alreay created connection number less than Max number of connection
// it will create new one. Otherwise it wil wait someone put connection back.
func (p *ConnPool) Get() *net.Conn {
	if p.isClosed() == true {
		return nil
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			fmt.Println(err)
			return
		}
		p.connChan <- conn
	}()
	select {
	case conn := <-p.connChan:
		return &conn
	default:
		return nil
	}

}

// Close close the connection pool. When close the connection pool it also close all connection already in connection pool.
// If connection not put back in connection it will not close. But it will close when it put back.
func (p *ConnPool) Close() {
	if p.isClosed() == true {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.closed = true
	close(p.connChan)
	for conn := range p.connChan {
		conn.Close()
	}
	return
}

// Put can put connection back in connection pool. If connection has been closed, the conneciton will be close too.
func (p *ConnPool) Put(conn net.Conn) {
	if p.isClosed() == true {
		return
	}
	if conn == nil {
		p.lock.Lock()
		p.totalConnNum = p.totalConnNum - 1
		p.lock.Unlock()
		return
	}

	select {
	case p.connChan <- conn:
		return
	default:
		conn.Close()
	}
	return
}

func (p *ConnPool) isClosed() bool {
	p.lock.Lock()
	ret := p.closed
	p.lock.Unlock()
	return ret
}

// createConn will create one connection from connCreator. And increase connection counter.
func (p *ConnPool) createConn() (net.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.totalConnNum >= p.maxConnNum {
		return nil, errors.New("")
	}
	conn, err := p.connCreator()
	if err != nil {
		return nil, err
	}
	p.totalConnNum = p.totalConnNum + 1
	return conn, nil
}
