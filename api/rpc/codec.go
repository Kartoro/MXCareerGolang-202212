package rpc

type Codec interface {
	ReadHeader()
	ReadBody()
	Write()
}

//func NewGobCodec(conn io.ReadWriteCloser) Codec {
//	dec := gob.NewDecoder()
//	enc := gob.NewEncoder()
//}