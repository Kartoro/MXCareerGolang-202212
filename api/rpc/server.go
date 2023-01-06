package rpc

import (
	"errors"
	"log"
	"net"
	"reflect"
	"strings"
)

type RemoteService struct {
	serviceMap map[string]*service // registered services
}

// rpcClient.Call(service.method, arg, reply)
type request struct {
	h     *Header // service.method
	arg   reflect.Value
	reply reflect.Value
	mtype *methodType // method name
	svc   *service    // service name
}

func NewServer() *RemoteService {
	return &RemoteService{make(map[string]*service)}
}

func (s *RemoteService) listen(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("network error: ", err)
	}
	return l, nil
}

func (s *RemoteService) lookupRpcMethod(serviceMethod string) (svc *service, mType *methodType, err error) {
	// separate service name and method name
	pos := strings.LastIndex(serviceMethod, ".")
	if pos < 0 {
		err = errors.New("rpc server: service/method request err : " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:pos], serviceMethod[pos+1:]

	svc, ok := s.serviceMap[serviceName]
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	mType = svc.methods[methodName]
	return
}

func (s *RemoteService) readRequest(cc Codec) (*request, error) {
	var h Header
	err := cc.ReadHeader(&h)
	if err != nil {
		return nil, err
	}

	req := &request{h: &h}
	req.svc, req.mtype, err = s.lookupRpcMethod(h.ServiceMethod)
	if err != nil {
		return nil, err
	}

	// get arg and reply type by reflection
	req.arg = req.mtype.newArg()
	req.reply = req.mtype.newReply()

	argPointer := req.arg.Interface()
	if req.arg.Type().Kind() != reflect.Ptr {
		argPointer = req.arg.Addr().Interface()
	}
	if err = cc.ReadBody(argPointer); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
	return req, nil
}

var invalidRequest = struct{}{}

func (s *RemoteService) handleRequest(cc Codec, req *request) error {
	err := req.svc.call(req.mtype, req.arg, req.reply)
	return err
}

func (s *RemoteService) sendResponse(cc Codec, h *Header, body interface{}) {
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (s *RemoteService) serveConn(conn net.Conn) {
	defer conn.Close()
	if conn == nil {
		log.Println("conn is nil")
		return
	}
	for {
		cc := NewGobCodec(conn)
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				// cc.Close() // it's not possible to recover, so close the connection
				return
			}
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest)
			return
		}

		err = s.handleRequest(cc, req)
		if err != nil {
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest)
			return
		}
		s.sendResponse(cc, req.h, req.reply.Interface())
	}
}

func (s *RemoteService) ListenAndServe(addr string) error {
	listener, err := s.listen(addr)
	if err != nil {
		return err
	}
	for {
		conn, e := listener.Accept()
		if e != nil {
			conn.Close()
			err = e
			break
		}
		go s.serveConn(conn)
	}
	return err
}

func (s *RemoteService) Register(service interface{}) {
	svc := newService(service)
	s.serviceMap[svc.name] = svc
}
