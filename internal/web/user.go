/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:27:11
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/web/user.go
 * @Description: 用户接口
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package web

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gz4z2b/go-webook/conf"
	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/service"
)

// UserHandler 我准备在上面定义跟用户有关的路由
type UserHandler struct {
	svc                       service.UserService
	emailExpersion            *regexp.Regexp
	passwordExpersion         *regexp.Regexp
	birthdayRegexExpersion    *regexp.Regexp
	nickNameRegexExpersion    *regexp.Regexp
	descriptionRegexExpersion *regexp.Regexp
}

// UserHandler构造方法
func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		passwordRegexpPattern   = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&_])[A-Za-z\d@$!%*?&_]{8,72}$`
		emailRegextPattern      = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		nickNameRegexPattern    = `^[\u4E00-\u9FFFa-zA-Z0-9!@#$%^&*()_+=\-[\]{}|\\:;"'<>,.?/~]{1,64}$`
		birthdayRegexPattern    = `^\d{4}-\d{2}-\d{2}$`
		descriptionRegexPattern = `^[\u4E00-\u9FFFa-zA-Z0-9!@#$%^&*()_+=\-[\]{}|\\:;。，？！……「」【】；'<>,.?/~]{1,10240}$`
	)
	passwordExpersion := regexp.MustCompile(passwordRegexpPattern, regexp.None)
	emailExpersion := regexp.MustCompile(emailRegextPattern, regexp.None)
	nickNameRegexExpersion := regexp.MustCompile(nickNameRegexPattern, regexp.None)
	birthdayRegexExpersion := regexp.MustCompile(birthdayRegexPattern, regexp.None)
	descriptionRegexExpersion := regexp.MustCompile(descriptionRegexPattern, regexp.None)

	return &UserHandler{
		emailExpersion:            emailExpersion,
		passwordExpersion:         passwordExpersion,
		nickNameRegexExpersion:    nickNameRegexExpersion,
		birthdayRegexExpersion:    birthdayRegexExpersion,
		descriptionRegexExpersion: descriptionRegexExpersion,
		svc:                       svc,
	}
}

// Signup 注册
func (u *UserHandler) Signup(ctx *gin.Context) {
	// 注册
	type signupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req signupReq

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

	// 用户存储
	err = u.svc.SignUp(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == service.ErrEmailConflict {
			ctx.String(http.StatusOK, "邮箱冲突啦~~")
			return
		}
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.String(http.StatusOK, "success")
}

// Login 登录
func (u *UserHandler) Login(ctx *gin.Context) {

	// 登录
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginReq
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.String(http.StatusOK, "输入数据格式错误")
		return
	}

	_, err = u.svc.Login(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrPasswordInvalid {
			ctx.String(http.StatusOK, "邮箱或密码错误")
			return
		}
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	userClaims := domain.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email:     req.Email,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	tokenStr, err := token.SignedString([]byte(conf.Keys.AuthorizationKey))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")
}

// Edit 修改
func (u *UserHandler) Edit(ctx *gin.Context) {
	// 修改
	type editReq struct {
		NickName    string `json:"nick_name"`
		BirthDay    string `json:"birth_day"`
		Description string `json:"description"`
	}
	var req editReq
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}

	ok, err := u.nickNameRegexExpersion.MatchString(req.NickName)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "昵称含非法字符")
		return
	}

	ok, err = u.birthdayRegexExpersion.MatchString(req.BirthDay)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}

	ok, err = u.descriptionRegexExpersion.MatchString(req.Description)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "简介含非法字符")
		return
	}

	email, exist := ctx.Get("user_email")
	if !exist {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		ctx.String(http.StatusOK, "登录态初始化错误")
		return
	}
	user, err := u.svc.FindByEmail(ctx, emailStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	birthDay, err := time.ParseInLocation("2006-01-02", req.BirthDay, time.Local)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	_, err = u.svc.AddProfile(ctx, user, &domain.Profile{
		NickName:    req.NickName,
		BirthDay:    birthDay.UnixMilli(),
		Description: req.Description,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "修改成功")

}

// Profile 个人档案
func (u *UserHandler) Profile(ctx *gin.Context) {
	// 个人档案
	email, exist := ctx.Get("user_email")
	if !exist {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		ctx.String(http.StatusOK, "登录态初始化错误")
		return
	}
	user, err := u.svc.FindByEmail(ctx, emailStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	profile, _ := u.svc.FindProfileByUser(ctx, user)

	ctx.JSON(http.StatusOK, profile)

}

// Logout 登出
func (u *UserHandler) Logout(ctx *gin.Context) {
	userClaims := domain.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email:     "",
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	tokenStr, err := token.SignedString([]byte(conf.Keys.AuthorizationKey))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "success")
}
