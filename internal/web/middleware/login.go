/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 20:45:58
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/web/middleware/login.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (loginMiddlewareBuilder *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	loginMiddlewareBuilder.paths = append(loginMiddlewareBuilder.paths, path)
	return loginMiddlewareBuilder
}

func (loginMiddlewareBuilder *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, value := range loginMiddlewareBuilder.paths {
			if value == ctx.Request.URL.Path {
				return
			}
		}

		session := sessions.Default(ctx)
		session.Options(sessions.Options{
			MaxAge: 30,
		})
		email := session.Get("user_email")
		if email == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user_email", email)
		now := time.Now().UnixMilli()
		lastLoginTime := session.Get("login_time")
		if lastLoginTime == nil {
			session.Set("login_time", now)
			if err := session.Save(); err != nil {
				panic(err)
			}
			return
		}
		lastLoginTimeVal, ok := lastLoginTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		if now-lastLoginTimeVal > 10000 {
			session.Set("login_time", now)
			session.Save()
		}

	}
}
