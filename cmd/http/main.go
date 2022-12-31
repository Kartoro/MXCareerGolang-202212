package main

import (
	"MXCareerGolang-202212/config"
	bff "MXCareerGolang-202212/internal/http"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("init http service...")
	bff.InitHttp()
}

func main() {
	// static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/images"))))

	http.HandleFunc("/login", bff.Login)
	http.HandleFunc("/index", bff.GetIndex)
	//http.HandleFunc("/profile", bff.Profile)
	//http.HandleFunc("/signup", bff.SignUp)
	//http.HandleFunc("/nickname", bff.UpdateNickname)
	//http.HandleFunc("/avatar", bff.UploadAvatar)
	// ...

	log.Println("http service started")
	if err := http.ListenAndServe(config.HTTPServerAddr, nil); err != nil {
		log.Println(err)
	}

}
