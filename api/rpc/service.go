package rpc

import (
	"go/ast"
	"log"
	"reflect"
)

type methodType struct {
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
}

type service struct {
	name    string
	t       reflect.Type
	v       reflect.Value
	methods map[string]*methodType
}

func (m *methodType) newArg() reflect.Value {
	var arg reflect.Value
	// arg may be a pointer type, or a value type
	if m.argType.Kind() == reflect.Ptr {
		arg = reflect.New(m.argType.Elem())
	} else {
		arg = reflect.New(m.argType).Elem()
	}
	return arg
}

func (m *methodType) newReply() reflect.Value {
	// reply must be a pointer type
	reply := reflect.New(m.replyType.Elem())
	switch m.replyType.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(m.replyType.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(m.replyType.Elem(), 0, 0))
	}
	return reply
}

func (s *service) call(m *methodType, arg reflect.Value, reply reflect.Value) error {
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.v, arg, reply})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
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

func newService(v interface{}) *service {
	s := new(service)
	s.v = reflect.ValueOf(v)
	s.name = reflect.Indirect(s.v).Type().Name()
	s.t = reflect.TypeOf(v)

	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	log.Println(s.name, s.v, s.t)
	s.registerMethods()
	return s
}
