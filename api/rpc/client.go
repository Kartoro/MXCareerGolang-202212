package rpc

import "MXCareerGolang-202212/internal/connection"

type RemoteClient struct {
	pool *connection.ConnPool
}

func ClientInit() (RemoteClient, error) {
	return RemoteClient{}, nil
}

func (r *RemoteClient) Call(serviceMethod string, arg interface{}, reply interface{}) (err error) {
	conn, err := r.pool.Get()
	defer r.pool.Put(conn)
	codec := NewGobCodec(*conn)
	codec.Write() // send encoded msg to RPC server
	codec.ReadHeader() // receive decoded msg from RPC sr
	codec.ReadBody()
	return nil
}
