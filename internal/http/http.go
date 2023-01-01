package http

import (
	"MXCareerGolang-202212/api/proto"
	"MXCareerGolang-202212/api/rpc"
	"MXCareerGolang-202212/config"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	loginTemplate   *template.Template
	profileTemplate *template.Template
	resultTemplate  *template.Template

	rpcClient *rpc.RemoteClient
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
	arg, err := loginReqConvRpcArg(rw, req)
	if err != nil {
		log.Println(err)
		return
	}
	resp := proto.LoginResponse{}
	if err = callLogin(rw, arg, &resp); err != nil {
		log.Println(err)
		return
	}
	handleLoginRet(rw, arg, &resp)
	log.Println("Login Done")

	return
}

func GetIndex(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		templateLogin(rw, TemplateLoginResponse{""})
	}
}

func loginReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.LoginRequest, error) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	if username == "" || password == "" {
		errMsg := "username or password is empty"
		templateLogin(rw, TemplateLoginResponse{errMsg})
		return &proto.LoginRequest{}, errors.New(errMsg)
	}
	arg := &proto.LoginRequest{
		Username: username,
		Password: password,
	}
	return arg, nil
}

// callLogin - RPC call
func callLogin(rw http.ResponseWriter, arg *proto.LoginRequest, resp *proto.LoginResponse) error {
	if err := rpcClient.Call("UserServices.Login", arg, resp); err != nil {
		log.Println(err)
		templateLogin(rw, TemplateLoginResponse{Msg: "RPC call failed"})
		return err
	}
	return nil
}

func handleLoginRet(rw http.ResponseWriter, arg *proto.LoginRequest, resp *proto.LoginResponse) {
	switch resp.Ret {
	case config.SUCCESS:
		// TODO: token
		templateLogin(rw, TemplateLoginResponse{Msg: "Login Success!"})
	default:
		templateLogin(rw, TemplateLoginResponse{Msg: "Login Failed!"})
	}
}