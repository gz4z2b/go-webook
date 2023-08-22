/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:44:46
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/web/init_web.go
 * @Description: 网络服务初始化
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package web

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
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
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	registerUserRoutes(server)

	return server
}

func registerUserRoutes(server *gin.Engine) {
	user := NewUserHandler()

	userGroup := server.Group("/users")
	userGroup.POST("/signup", user.Signup)
	userGroup.GET("/:uid", user.Profile)
	userGroup.POST("/login", user.Login)
	userGroup.POST("/edit", user.Edit)
}
