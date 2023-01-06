package http

import (
	"MXCareerGolang-202212/api/proto"
	"MXCareerGolang-202212/api/rpc"
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/model"
	"MXCareerGolang-202212/internal/util"
	"errors"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	loginTemplate   *template.Template
	profileTemplate *template.Template
	resultTemplate  *template.Template
	rpcClient       rpc.RemoteClient
)

type TemplateLoginResponse struct {
	Msg string
}

type TemplateProfileResponse struct {
	UserName string
	NickName string
	PicName  string
}

func InitHttp() {
	loginTemplate = template.Must(template.ParseFiles("./static/templates/login.html"))
	profileTemplate = template.Must(template.ParseFiles("./static/templates/profile.html"))
	resultTemplate = template.Must(template.ParseFiles("./static/templates/result.html"))

	rpcClient, _ = rpc.ClientInit()
}

func SignUp(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		arg, err := signUpReqConvRpcArg(rw, req)
		if err != nil {
			return
		}
		reply := proto.LoginResponse{}
		if err = callSignUp(rw, arg, &reply); err != nil {
			return
		}
		handleSignUpRet(rw, arg, &reply)
	}
}

func Login(rw http.ResponseWriter, req *http.Request) {
	arg, err := loginReqConvRpcArg(rw, req)
	if err != nil {
		return
	}
	reply := proto.LoginResponse{}
	if err = callLogin(rw, arg, &reply); err != nil {
		return
	}
	handleLoginRet(rw, arg, &reply)
}

func GetIndex(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		templateLogin(rw, TemplateLoginResponse{Msg: ""})
	}
}

func GetProfile(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		arg, err := getProfileReqConvRpcArg(rw, req)
		if err != nil {
			return
		}
		reply := proto.ProfileResponse{}
		if err = callGetProfile(rw, arg, &reply); err != nil {
			return
		}
		handleGetProfileRet(rw, arg, &reply)
	}
}

func UpdateNickName(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		arg, err := updateNickNameReqConvRpcArg(rw, req)
		if err != nil {
			return
		}
		reply := proto.ProfileResponse{}
		if err = callUpdateNickName(rw, arg, &reply); err != nil {
			return
		}
		handleUpdateNickNameRet(rw, arg, &reply)
	}
}

func UploadProfilePicture(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		arg, err := uploadProfilePictureReqConvRpcArg(rw, req)
		if err != nil {
			return
		}
		reply := proto.ProfileResponse{}
		if err = callUploadProfilePicture(rw, arg, &reply); err != nil {
			return
		}
		handleUploadProfilePictureRet(rw, arg, &reply)
	}
}

// Call RPC method: UserServices.SignUp
func callSignUp(rw http.ResponseWriter, arg *proto.LoginRequest, reply *proto.LoginResponse) error {
	if err := rpcClient.Call("UserServices.SignUp", arg, reply); err != nil {
		log.Printf("Call signUp failed. username:%s, err:%q\n", arg.UserName, err)
		_, _ = rw.Write([]byte("Sign Up failed!"))
		return err
	}
	return nil
}

func handleSignUpRet(rw http.ResponseWriter, arg *proto.LoginRequest, reply *proto.LoginResponse) {
	switch reply.Ret {
	case model.Success:
		_, _ = rw.Write([]byte("Sign Up succeed!"))

	case model.EmptyUserNameOrPwdError:
		_, _ = rw.Write([]byte("Username or password could not be empty!"))

	default:
		_, _ = rw.Write([]byte("Sign Up failed!"))
	}
	log.Printf("handleSignUpRet done. username:%s, ret:%d\n", arg.UserName, reply.Ret)
}

func signUpReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.LoginRequest, error) {
	userName := req.FormValue("username")
	password := req.FormValue("password")
	nickName := req.FormValue("nickname")

	if userName == "" || password == "" {
		_, _ = rw.Write([]byte("username and password couldn't be NULL!"))
		return &proto.LoginRequest{}, errors.New("username and password couldn't be NULL")
	}
	log.Printf("username = %s, password = %s, nickname = %s\n", userName, password, nickName)
	arg := proto.LoginRequest{
		UserName: userName,
		Password: password,
		NickName: nickName,
	}
	return &arg, nil
}

// Call RPC method: UserServices.Login
func callLogin(rw http.ResponseWriter, arg *proto.LoginRequest, reply *proto.LoginResponse) error {
	if err := rpcClient.Call("UserServices.Login", arg, &reply); err != nil {
		log.Printf("Call login failed. username:%s, err:%q\n", arg.UserName, err)
		templateLogin(rw, TemplateLoginResponse{Msg: "Login failed!"})
		return err
	}
	return nil
}

func handleLoginRet(rw http.ResponseWriter, arg *proto.LoginRequest, reply *proto.LoginResponse) {
	switch reply.Ret {
	case model.Success: // send username & token back to client by cookie
		cookie := http.Cookie{Name: "username", Value: arg.UserName, MaxAge: config.MaxExTime,
			HttpOnly: true, SameSite: http.SameSiteStrictMode}
		http.SetCookie(rw, &cookie)
		cookie = http.Cookie{Name: "token", Value: reply.Token, MaxAge: config.MaxExTime,
			HttpOnly: true, SameSite: http.SameSiteStrictMode}
		http.SetCookie(rw, &cookie)
		templateResult(rw, TemplateLoginResponse{Msg: "Login succeed!"})

	case model.UserNameOrPasswordError:
		templateLogin(rw, TemplateLoginResponse{Msg: "Wrong username or password!"})

	default:
		templateLogin(rw, TemplateLoginResponse{Msg: "Login failed!"})
	}
	//log.Printf("login done. username:%s, ret:%d\n", arg.UserName, reply.Ret)
}

func loginReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.LoginRequest, error) {
	userName := req.FormValue("username")
	password := req.FormValue("password")
	if userName == "" || password == "" {
		templateLogin(rw, TemplateLoginResponse{Msg: "Username and password couldn't be empty!"})
		return &proto.LoginRequest{}, errors.New("username and password couldn't be empty")
	}

	arg := proto.LoginRequest{
		UserName: userName,
		Password: password,
	}

	return &arg, nil
}

// Call RPC method: UserServices.GetProfile
func callGetProfile(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) error {
	if err := rpcClient.Call("UserServices.GetProfile", arg, &reply); err != nil {
		log.Printf("Call getProfile failed. username:%s, err:%q\n", arg.UserName, err)
		templateLogin(rw, TemplateLoginResponse{Msg: "Session expired, please login again"})
		return err
	}
	return nil
}

func handleGetProfileRet(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) {
	switch reply.Ret {
	case model.Success:
		if reply.PicName == "" {
			reply.PicName = config.DefaultImagePath
		}
		log.Println(reply)
		templateProfile(rw, TemplateProfileResponse{
			UserName: reply.UserName,
			NickName: reply.NickName,
			PicName:  reply.PicName})

	case model.TokenUnmatchedError:
		templateLogin(rw, TemplateLoginResponse{Msg: "Session expired, please login again"})

	case model.NilDataError:
		templateResult(rw, TemplateLoginResponse{Msg: "User not exists!"})

	default:
		templateResult(rw, TemplateLoginResponse{Msg: "Failed to get user profile!"})
	}
	log.Printf("GetProfile done. username:%s, ret:%d\n", arg.UserName, reply.Ret)
}

func getProfileReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.ProfileRequest, error) {
	token, err := getToken(rw, req)
	if err != nil {
		return &proto.ProfileRequest{}, err
	}
	userName := req.FormValue("username")
	if userName == "" {
		nameCookie, err := req.Cookie("username")
		if err != nil {
			templateLogin(rw, TemplateLoginResponse{Msg: ""}) // if usernames in both form and cookie are empty, login again
			return &proto.ProfileRequest{}, errors.New("username empty")
		}
		userName = nameCookie.Value
	}

	arg := proto.ProfileRequest{
		UserName: userName,
		Token:    token,
	}

	return &arg, nil
}

// Call RPC method: UserServices.UpdateNickName
func callUpdateNickName(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) error {
	if err := rpcClient.Call("UserServices.UpdateNickName", arg, &reply); err != nil {
		log.Printf("Call updateNickName failed. username:%s, err:%q", arg.UserName, err)
		templateResult(rw, TemplateLoginResponse{Msg: "Update nickname failed!"})
		return err
	}
	return nil
}

func handleUpdateNickNameRet(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) {
	switch reply.Ret {
	case model.Success:
		templateResult(rw, TemplateLoginResponse{Msg: "Update nickname succeed!"})

	case model.TokenUnmatchedError:
		templateLogin(rw, TemplateLoginResponse{Msg: "Session expired, please login again"})

	case model.UserNotExistError:
		templateResult(rw, TemplateLoginResponse{Msg: "User not exists!"})

	default:
		templateResult(rw, TemplateLoginResponse{Msg: "Update nickname failed!"})

	}
	log.Printf("update nickname done. username:%s, nickname:%s, ret:%d", arg.UserName, arg.NickName, reply.Ret)
}

func updateNickNameReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.ProfileRequest, error) {
	token, err := getToken(rw, req)
	if err != nil {
		return &proto.ProfileRequest{}, err
	}
	userName := req.FormValue("username")
	nickName := req.FormValue("nickname")

	arg := proto.ProfileRequest{
		UserName: userName,
		NickName: nickName,
		Token:    token,
	}
	return &arg, nil
}

// Call RPC method: UserServices.UpdateAvatar
func callUploadProfilePicture(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) error {
	if err := rpcClient.Call("UserServices.UpdateAvatar", arg, &reply); err != nil {
		log.Printf("Call uploadProfilePic failed. username:%s, err:%q", arg.UserName, err)
		templateResult(rw, TemplateLoginResponse{Msg: "Avatar updated failed!"})
		return err
	}
	return nil
}

func handleUploadProfilePictureRet(rw http.ResponseWriter, arg *proto.ProfileRequest, reply *proto.ProfileResponse) {
	switch reply.Ret {
	case model.Success:
		templateResult(rw, TemplateLoginResponse{Msg: "Avatar updated!"})

	case model.TokenUnmatchedError:
		templateLogin(rw, TemplateLoginResponse{Msg: "Session expired, please login again."})

	case model.UserNotExistError:
		templateResult(rw, TemplateLoginResponse{Msg: "User not exist!"})

	default:
		templateResult(rw, TemplateLoginResponse{Msg: "Avatar updated failed!"})
	}
	log.Printf("Upload picture done. username:%s, filepath:%s, ret:%d", arg.UserName, arg.FileName, reply.Ret)
}

func uploadProfilePictureReqConvRpcArg(rw http.ResponseWriter, req *http.Request) (*proto.ProfileRequest, error) {
	token, err := getToken(rw, req)
	if err != nil {
		return &proto.ProfileRequest{}, err
	}
	userName := req.FormValue("username")
	fileName, err := getFileName(userName, rw, req)
	if err != nil {
		return &proto.ProfileRequest{}, err
	}

	arg := proto.ProfileRequest{
		UserName: userName,
		FileName: fileName,
		Token:    token,
	}
	return &arg, nil
}

func getToken(rw http.ResponseWriter, req *http.Request) (string, error) {
	token, err := req.Cookie("token")
	if err != nil {
		log.Printf("get token failed. err:%q", err)
		templateLogin(rw, TemplateLoginResponse{})
		return "", errors.New("get token failed")
	}

	return token.Value, nil
}

func getFileName(userName string, rw http.ResponseWriter, req *http.Request) (string, error) {
	file, head, err := req.FormFile("image")
	if err != nil {
		templateResult(rw, TemplateLoginResponse{Msg: "get file failed！"})
		log.Printf("get file name failed. username:%s, err:%q", userName, err)
		return "", errors.New("get file name failed")
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// generate file path
	newName, isLegal := util.CheckAndCreateFileName(head.Filename)
	if !isLegal {
		templateResult(rw, TemplateLoginResponse{Msg: "file type illegal! Must be .jpg/.jpeg/.png/.gif"})
		return "", errors.New("file type illegal")
	}
	filePath := config.StaticFilePath + "/" + newName
	fileName := newName
	log.Println(filePath)

	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		templateResult(rw, TemplateLoginResponse{Msg: "file open failed！"})
		return "", err
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {

		}
	}(dstFile)

	// copy and upload avatar to static folder
	_, err = io.Copy(dstFile, file)
	if err != nil {
		templateResult(rw, TemplateLoginResponse{Msg: "file upload failed!！"})
		return "", err
	}

	return fileName, nil
}

func templateLogin(rw http.ResponseWriter, reply TemplateLoginResponse) {
	if err := loginTemplate.Execute(rw, reply); err != nil {
		log.Printf("template login failed: %q\n", err)
	}
}

func templateProfile(rw http.ResponseWriter, reply TemplateProfileResponse) {
	if err := profileTemplate.Execute(rw, reply); err != nil {
		log.Printf("template profile failed: %q\n", err)
	}
}

func templateResult(rw http.ResponseWriter, reply TemplateLoginResponse) {
	if err := resultTemplate.Execute(rw, reply); err != nil {
		log.Printf("template result failed: %q\n", err)
	}
}
