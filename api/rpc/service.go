package rpc

import (
	"go/ast"
	"log"
	"reflect"
)

type UserServices struct{}

type service struct {
	name    string
	t       reflect.Type
	v       reflect.Value
	methods map[string]*methodType
}

type methodType struct {
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
}

func newService(v interface{}) *service {
	s := new(service)
	s.v = reflect.ValueOf(v)
	s.name = reflect.Indirect(s.v).Type().Name()
	s.t = reflect.TypeOf(v)

	if !ast.IsExported(s.name) {
		log.Fatalf("newService error")
	}
	s.registerMethods()
	log.Println(s.name, s.v, s.t, " registered")
	return s
}

func (s *service) registerMethods() {
	s.methods = make(map[string]*methodType)
	for i := 0; i < s.t.NumMethod(); i++ {
		method := s.t.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}

		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		s.methods[method.Name] = &methodType{
			method:    method,
			argType:   argType,
			replyType: replyType,
		}
		log.Printf("RPC registerMethods: %s.%s\n", s.name, method.Name)
	}
}
