package http

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	loginTemplate   *template.Template
	profileTemplate *template.Template
	resultTemplate  *template.Template

)

type TemplateLoginResponse struct {
	Msg string
}

type TemplateProfileResponse struct {
	Username string
	Nickname string
	Avatar   string
}

type TemplateResultResponse struct {
	Msg string
}

func InitHttp() {
	loginTemplate = template.Must(template.ParseFiles("./static/templates/login.html"))
	profileTemplate = template.Must(template.ParseFiles("./static/templates/profile.html"))
	resultTemplate = template.Must(template.ParseFiles("./static/templates/result.html"))
}

func templateLogin(rw http.ResponseWriter, reply TemplateLoginResponse) {
	if err := loginTemplate.Execute(rw, reply); err != nil {
		fmt.Println(err)
	}
}

func templateProfile(rw http.ResponseWriter, reply TemplateProfileResponse) {
	if err := profileTemplate.Execute(rw, reply); err != nil {
		fmt.Println(err)
	}
}

func templateResult(rw http.ResponseWriter, reply TemplateResultResponse) {
	if err := resultTemplate.Execute(rw, reply); err != nil {
		fmt.Println(err)
	}
}

func Login(rw http.ResponseWriter, req *http.Request) {
	// TODO: rpc
	fmt.Println("Login")
	return
}

func GetIndex(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		templateLogin(rw, TemplateLoginResponse{""})
	}
}

