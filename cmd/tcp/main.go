package main

import (
	"MXCareerGolang-202212/api/rpc"
	"MXCareerGolang-202212/config"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("tcp service started")
	var services rpc.UserServices
	s := rpc.NewServer()
	s.Register(&services)
	if err := http.ListenAndServe(config.TCPServerAddr, nil); err != nil {
		log.Println(err)
	}
}
