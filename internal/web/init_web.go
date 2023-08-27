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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gz4z2b/go-webook/internal/repository"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	"github.com/gz4z2b/go-webook/internal/service"
	"github.com/gz4z2b/go-webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("WnXqVdx2tPCKRQMPb2L6YXoKaHvbDvMd"), []byte("NIVL94TXGmuKmFON6ud71dulRCjXALMy"))
	store, err := redis.NewStore(16, "tcp", "localhost:13317", "", []byte("WnXqVdx2tPCKRQMPb2L6YXoKaHvbDvMd"), []byte("NIVL94TXGmuKmFON6ud71dulRCjXALMy"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("login", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/signup").IgnorePath("/users/login").Build())

	registerUserRoutes(server)

	return server
}

func registerUserRoutes(server *gin.Engine) {
	user := initUser()

	userGroup := server.Group("/users")
	userGroup.POST("/signup", user.Signup)
	userGroup.GET("/:uid", user.Profile)
	userGroup.POST("/login", user.Login)
	userGroup.POST("/edit", user.Edit)
}

func initUser() *UserHandler {
	db, err := gorm.Open(mysql.Open("root:gz4z2b@tcp(127.0.0.1:13316)/webook"), &gorm.Config{})
	if err != nil {
		panic("数据库初始化失败")
	}
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	user := NewUserHandler(userService)
	return user
}
