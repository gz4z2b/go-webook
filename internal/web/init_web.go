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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gz4z2b/go-webook/conf"
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
			return strings.Contains(origin, "webook.gdtengnan.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("WnXqVdx2tPCKRQMPb2L6YXoKaHvbDvMd"), []byte("NIVL94TXGmuKmFON6ud71dulRCjXALMy"))
	//store, err := redis.NewStore(16, "tcp", fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port), conf.Redis.Password, []byte(conf.Keys.AuthorizationKey), []byte(conf.Keys.EncryptKey))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("login", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/signup").IgnorePath("/users/login").IgnorePath("/hello").Build())

	registerUserRoutes(server)

	return server
}

func registerUserRoutes(server *gin.Engine) {
	user := initUser()

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

func initUser() *UserHandler {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Db.User, conf.Db.Password, conf.Db.Host, conf.Db.Port, conf.Db.Db)), &gorm.Config{})
	if err != nil {
		panic("数据库初始化失败")
	}
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	user := NewUserHandler(userService)
	return user
}
