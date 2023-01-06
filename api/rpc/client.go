package rpc

import (
	pool "MXCareerGolang-202212/internal/connection"
	"errors"
	"io"
	"log"
	"time"
)

type RemoteClient struct {
	pool    *pool.ConnPool
	timeout time.Duration
	closed  bool
}

// ClientInit - init client with connection pool
// n - connection pool size; address - tcp server address
func ClientInit() (RemoteClient, error) {
	client := RemoteClient{}
	client.pool = pool.InitialisePoolValue(client.pool)
	log.Println("connection pool successfully initialized!")
	return client, nil
}

// Call - RPC RemoteClient invoke this method to get reply from RemoteService
// serviceMethod - RPC method name; arg - RPC method arguments; reply - store reply from RemoteService (must be pointer type)
func (r *RemoteClient) Call(serviceMethod string, arg interface{}, reply interface{}) (err error) {
	conn := r.pool.Get()
	defer r.pool.Put(*conn)
	cc := NewGobCodec(*conn)
	var h = &Header{
		ServiceMethod: serviceMethod,
		Error:         "",
	}

	// send to RPC server
	if err = cc.Write(h, arg); err != nil {
		return err
	}
	// receive from RPC server
	if err = cc.ReadHeader(h); err != nil {
		return err
	}
	if h.Error != "" {
		return errors.New(h.Error)
	}
	err = cc.ReadBody(reply)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
