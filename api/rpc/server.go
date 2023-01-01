package rpc

type RemoteService struct {
	serviceMap map[string]*service
}

func NewServer() *RemoteService{
	return &RemoteService{make(map[string]*service)}
}

func (s *RemoteService) Register(service interface{}) {
	svc := newService(service)
	s.serviceMap[svc.name] = svc
}

