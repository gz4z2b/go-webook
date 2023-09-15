// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/gz4z2b/go-webook/internal/repository"
	"github.com/gz4z2b/go-webook/internal/repository/cache"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	"github.com/gz4z2b/go-webook/internal/service"
	"github.com/gz4z2b/go-webook/internal/web"
)

// Injectors from wire.go:

func InitWebService() *gin.Engine {
	db := InitDb()
	userDAO := dao.NewUserDAO(db)
	cmdable := InitCache()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	v := web.InitUserMidleware()
	engine := web.InitWebService(userHandler, v)
	return engine
}