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

import "github.com/gin-gonic/gin"

func RegisterRoutes() *gin.Engine {
	server := gin.Default()

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
