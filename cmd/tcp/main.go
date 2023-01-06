package main

import (
	"MXCareerGolang-202212/api/rpc"
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/tcp/service"
	"log"
	"net/http"
)

func main() {
	go func() {
		err := http.ListenAndServe(config.PprofAddr, nil)
		if err != nil {
			return
		}
	}()

	var services service.UserServices
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := rpc.NewServer()
	s.Register(&services)
	err := s.ListenAndServe(config.TCPServerAddr)
	if err != nil {
		log.Fatalf("ListenAndServe failed! err: %v\n", err)
		return
	}
}
