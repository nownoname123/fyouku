package controllers

import (
	"context"
	userRpcProto "fyoukuapi/micro/user/proto"
	"fyoukuapi/model"
)

func UserRpcController(ctx context.Context, req *userRpcProto.RequestLogin, res *userRpcProto.ResponseLogin) error {
	var (
		userLoginProto userRpcProto.LoginUser

		uid   int64
		uname string
		err   error
	)
	var user UserType
	user.Mobile = req.Mobile
	user.Password = req.Password

	if user.Mobile == "" {
		res.Code = 4001
		res.Msg = "手机号不能为空"
		goto ERR
	}

	if user.Password == "" {
		res.Code = 4003
		res.Msg = "密码不能为空"
		goto ERR
	}
	uid, uname, err = model.UserLogin(user.Mobile, user.Password)
	if uid != 0 {
		userLoginProto.Uid = uid
		userLoginProto.Username = uname
		res.Code = 0
		res.Msg = "登陆成功"
		res.Items = &userLoginProto
		res.Count = 1
		return nil

	}
	if err == nil {
		res.Code = 4004
		res.Msg = "手机号或密码不正确"
		goto ERR
	}

ERR:
	res.Items = &userLoginProto
	res.Count = 0
	return nil
}
