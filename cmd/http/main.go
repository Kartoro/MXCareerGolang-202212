package main

import (
	"MXCareerGolang-202212/config"
	httpService "MXCareerGolang-202212/internal/http"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// init - parse static html templates
func init() {
	httpService.InitHttp()
}

func main() {
	//go func() {
	//	http.ListenAndServe(config.PprofAddr, nil)
	//}()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init rpc connection
	var err error
	if err != nil {
		log.Fatal("Failed to connect with rpc server!")
		return
	}

	// static file server. NOTE: for Handle and StripPrefix, the format must be /<path>/ !!!
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/images"))))

	http.HandleFunc("/index", httpService.GetIndex)
	http.HandleFunc("/profile", httpService.GetProfile)
	http.HandleFunc("/login", httpService.Login)
	http.HandleFunc("/signup", httpService.SignUp)
	http.HandleFunc("/nickname", httpService.UpdateNickName)
	http.HandleFunc("/avatar", httpService.UploadProfilePicture)

	err = http.ListenAndServe(config.HTTPServerAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
