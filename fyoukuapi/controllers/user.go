package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
)

type UserType struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

func UserRegister(c *gin.Context) {
	var newUser UserType
	//	err := c.BindJSON(&newUser)
	newUser.Mobile = c.PostForm("mobile")
	newUser.Password = c.PostForm("password")
	var err error
	err = nil

	if newUser.Mobile == "" {
		ReturnError(c, 4002, "手机号不能为空")
		return
	}
	if newUser.Password == "" {
		ReturnError(c, 4003, "密码不能为空")
		return
	}
	//传递数据给model，去数据库中访问看看手机号是否被注册
	pd, err := model.IsUserMobile(newUser.Mobile)
	if err != nil {
		ReturnError(c, 4005, "访问数据库错误")
		return
	}
	if !pd {
		ReturnError(c, 4004, "此手机号已被注册")
		return
	}
	err = model.UserSave(newUser.Mobile, newUser.Password)
	if err != nil {
		ReturnError(c, 4006, "注册失败请重试")
		return
	}
	ReturnSuccess(c, 0, "注册成功", nil, 0)
}
func UserLogin(c *gin.Context) {
	var user UserType
	user.Mobile = c.PostForm("mobile")
	user.Password = c.PostForm("password")

	if user.Mobile == "" {
		ReturnError(c, 4002, "手机号不能为空")
		return
	}
	if user.Password == "" {
		ReturnError(c, 4003, "密码不能为空")
		return
	}
	uid, uname, err := model.UserLogin(user.Mobile, user.Password)
	if uid != 0 {
		ReturnSuccess(c, 0, "登陆成功", map[string]interface{}{"uid": uid, "username": uname}, 1)
		return
	}
	if err == nil {
		ReturnError(c, 4007, "手机号或者密码错误")
		return
	}
	ReturnError(c, 4005, "访问数据库错误")
}
