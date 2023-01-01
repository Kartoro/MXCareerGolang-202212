package connection

import (
	"errors"
	"net"
	"sync"
)

type Pool interface {
	Get() (*net.Conn, error)
	Put(*net.Conn) error
	Close()
	Remove(conn *net.Conn) error
}

type ConnPool struct {
	lock sync.Mutex
	connChan chan *net.Conn
	totalConnNum, maxConnNum int32
	connCreator func() (*net.Conn, error)
}

func (p *ConnPool) Get() (*net.Conn, error){
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.connChan <- conn
	}()
	select {
	case conn := <- p.connChan:
		return conn, nil
	default:
		return nil, nil
	}
}

func (p *ConnPool) Put(conn *net.Conn) error {
	return nil
}

func (p *ConnPool) createConn() (*net.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.totalConnNum >= p.maxConnNum {
		return nil, errors.New("max conn exceed")
	}
	conn, err := p.connCreator()
	if err != nil {
		return nil, err
	}
	p.totalConnNum += 1
	return conn, nil
}