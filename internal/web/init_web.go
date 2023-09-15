/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:44:46
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/web/init_web.go
 * @Description: 网络服务初始化
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gz4z2b/go-webook/internal/web/middleware"
)

func InitWebService(userHandler *UserHandler, mids []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(mids...)
	registerUserRoutes(server, userHandler)
	return server
}

func InitUserMidleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins: []string{"*"},
			//AllowMethods: []string{"POST", "GET"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 你不加这个，前端是拿不到的
			ExposeHeaders: []string{"x-jwt-token"},
			// 是否允许你带 cookie 之类的东西
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					// 你的开发环境
					return true
				}
				return strings.Contains(origin, "webook.gdtengnan.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/signup").IgnorePath("/users/login").IgnorePath("/hello").Build(),
	}
}

func registerUserRoutes(server *gin.Engine, user *UserHandler) {

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	})

	userGroup := server.Group("/users")
	userGroup.POST("/signup", user.Signup)
	userGroup.GET("/:uid", user.Profile)
	userGroup.POST("/login", user.Login)
	userGroup.POST("/edit", user.Edit)
	userGroup.POST("/logout", user.Logout)
}
