/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:27:11
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/web/user.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package web

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

// UserHandler 我准备在上面定义跟用户有关的路由
type UserHandler struct {
	emailExpersion    *regexp.Regexp
	passwordExpersion *regexp.Regexp
}

// UserHandler构造方法
func NewUserHandler() *UserHandler {
	const (
		passwordRegexpPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&_])[A-Za-z\d@$!%*?&_]{8,}$`
		emailRegextPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	)
	passwordExpersion := regexp.MustCompile(passwordRegexpPattern, regexp.None)
	emailExpersion := regexp.MustCompile(emailRegextPattern, regexp.None)

	return &UserHandler{
		emailExpersion:    emailExpersion,
		passwordExpersion: passwordExpersion,
	}
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	// 注册
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq

	if err := ctx.BindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "数据格式错误")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err := u.passwordExpersion.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码复杂度不够")
		return
	}

	ok, err = u.emailExpersion.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	ctx.String(http.StatusOK, "success")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	// 登录
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	// 修改
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	// 个人档案
}
