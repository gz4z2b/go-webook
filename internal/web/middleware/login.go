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
		email := session.Get("user_email")
		if email == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user_email", email)

		return
	}
}
