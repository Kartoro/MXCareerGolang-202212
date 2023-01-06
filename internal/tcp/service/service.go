package service

import (
	"MXCareerGolang-202212/api/proto"
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/model"
	"MXCareerGolang-202212/internal/tcp/dao"
	"MXCareerGolang-202212/internal/tcp/redis"
	"MXCareerGolang-202212/internal/util"
	"log"
)

type UserServices struct{}

func (s *UserServices) SignUp(req *proto.LoginRequest, resp *proto.LoginResponse) error {
	if req.UserName == "" || req.Password == "" {
		resp.Ret = model.EmptyUserNameOrPwdError
		return nil
	}
	if req.NickName == "" {
		req.NickName = req.UserName
	}

	if err := dao.CreateUser(req.UserName, req.NickName, req.Password); err != nil {
		resp.Ret = model.UserNameDupError
		log.Printf("create user failed. usernam:%s, err:%q\n", req.UserName, err)
		return nil
	}

	resp.Ret = model.Success
	return nil
}

func (s *UserServices) Login(req *proto.LoginRequest, resp *proto.LoginResponse) error {
	// check authentication
	ok, err := loginAuth(req)
	if err != nil {
		resp.Ret = model.UserLoginError
		return nil
	}
	if !ok {
		resp.Ret = model.UserNameOrPasswordError
		return nil
	}

	// generate token
	token := util.GetToken(req.UserName)
	err = updateCache(req, token)
	if err != nil {
		resp.Ret = model.UserLoginError
		return nil
	}

	// login successfully
	resp.Ret = model.Success
	resp.Token = token
	//log.Printf("username:%s login successfully.\n", req.UserName)
	return nil
}

func updateCache(req *proto.LoginRequest, token string) error {
	redis.SetPassword(req.UserName, req.Password)
	err := redis.SetToken(req.UserName, token, int64(config.MaxExTime))
	if err != nil {
		log.Printf("updateCache SetToken failed. username:%s, token:%s, err:%q\n", req.UserName, token, err)
	}
	return nil
}

func loginAuth(req *proto.LoginRequest) (bool, error) {
	// check redis first
	ok := redis.LoginAuth(req.UserName, req.Password)
	if ok {
		//log.Printf("login through redis >> username:%s\n", req.UserName)
		return true, nil
	}

	// if check redis failed (no data or err), check mysql
	ok, err := dao.LoginAuth(req.UserName, req.Password)
	if err != nil {
		log.Printf("login failed. username:%s, password:%s, err:%q\n", req.UserName, req.Password, err)
		return false, nil
	}
	//log.Printf("login through mysql >> username:%s\n", req.UserName)
	return ok, nil
}

func (s *UserServices) GetProfile(req *proto.ProfileRequest, resp *proto.ProfileResponse) error {
	// 校验token
	ok, err := redis.CheckToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = model.GetUserProfileError
		log.Printf("check token failed. username:%s, token:%s, err:%q\n", req.UserName, req.Token, err)
		return nil
	}
	if !ok {
		resp.Ret = model.TokenUnmatchedError
		return nil
	}

	userProfile, isRead := getUserProfile(req)
	if !isRead {
		resp.Ret = model.NilDataError
		return nil
	}

	log.Printf("get profile done. username:%s\n", req.UserName)
	*resp = proto.ProfileResponse{Ret: model.Success, UserName: req.UserName, NickName: userProfile.NickName, PicName: userProfile.Avatar}
	return nil
}

func getUserProfile(req *proto.ProfileRequest) (model.User, bool) {
	// check redis first
	userProfile, isRead := redis.GetProfile(req.UserName)
	if isRead {
		log.Printf("redis get profile done. username:%s\n", req.UserName)
		return userProfile, true
	}

	// if check redis failed (no data or err), check mysql
	userProfile, isRead = dao.GetProfile(req.UserName)
	if !isRead {
		log.Printf("dao get profile failed. username:%s\n", req.UserName)
		return model.User{}, false
	}

	// update redis
	err := redis.SetProfile(req.UserName, userProfile.NickName, userProfile.Avatar)
	if err != nil {
		log.Printf("dao get profile failed. username:%s\n", req.UserName)
	}

	return userProfile, true
}

func (s *UserServices) UpdateAvatar(req *proto.ProfileRequest, resp *proto.ProfileResponse) error {
	ok, err := redis.CheckToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("check token failed. username:%s, token:%s, err:%q\n", req.UserName, req.Token, err)
		return nil
	}
	if !ok {
		resp.Ret = model.TokenUnmatchedError
		return nil
	}

	// delete cache for consistency
	if err := redis.DelCache(req.UserName); err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("DelCache failed. username:%s, err:%q\n", req.UserName, err)
		return nil
	}

	ok, err = dao.UpdateAvatar(req.UserName, req.FileName)
	if err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("UpdateAvatar failed. username:%s, filename:%s, err:%q\n", req.UserName, req.FileName, err)
		return nil
	}
	if !ok {
		resp.Ret = model.UserNotExistError
		return nil
	}
	resp.Ret = model.Success
	log.Printf("UpdateAvatar done. username:%s, filename:%s\n", req.UserName, req.FileName)
	return nil
}

func (s *UserServices) UpdateNickName(req *proto.ProfileRequest, resp *proto.ProfileResponse) error {
	ok, err := redis.CheckToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("check token failed. username:%s, token:%s, err:%q\n", req.UserName, req.Token, err)
		return nil
	}
	if !ok {
		resp.Ret = model.TokenUnmatchedError
		return nil
	}

	// delete cache for consistency
	if err := redis.DelCache(req.UserName); err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("redis invaild cache failed. username:%s, err:%q\n", req.UserName, err)
		return nil
	}

	ok, err = dao.UpdateNickName(req.UserName, req.NickName)
	if err != nil {
		resp.Ret = model.UserUpdateError
		log.Printf("UpdateNickName failed. username:%s, nickname:%s, err:%q\n", req.UserName, req.NickName, err)
		return nil
	}
	if !ok {
		resp.Ret = model.UserNotExistError
		return nil
	}
	resp.Ret = model.Success
	log.Printf("UpdateNickName done. username:%s, nickname:%s\n", req.UserName, req.NickName)
	return nil
}
