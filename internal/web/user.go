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

import "github.com/gin-gonic/gin"

// UserHandler 我准备在上面定义跟用户有关的路由
type UserHandler struct {
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	// 注册
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq
	err := ctx.Bind(req)
	if err != nil {
		return
	}
	ctx.JSON(200, gin.H{
		"message": "success",
	})
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
