package rpc

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type Header struct {
	ServiceMethod string
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(h *Header, body interface{}) error
}

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func (cc *GobCodec) Close() error {
	return cc.conn.Close()
}

func (cc *GobCodec) ReadHeader(h *Header) error {
	return cc.dec.Decode(h)
}

func (cc *GobCodec) ReadBody(body interface{}) error {
	return cc.dec.Decode(body)
}

func (cc *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = cc.buf.Flush()
		if err != nil {
			_ = cc.Close()
		}
	}()

	if err := cc.enc.Encode(h); err != nil {
		log.Println("gob encoding header error", err)
		return err
	}

	if err = cc.enc.Encode(body); err != nil {
		log.Println("gob encoding body error", err)
		return err
	}

	return nil
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
