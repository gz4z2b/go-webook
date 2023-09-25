//go:build wireinject

/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 11:06:03
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/ioc/wire.go
 * @Description: 初始化服务
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gz4z2b/go-webook/internal/repository"
	"github.com/gz4z2b/go-webook/internal/repository/cache"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	"github.com/gz4z2b/go-webook/internal/service"
	"github.com/gz4z2b/go-webook/internal/web"
)

func InitWebService() *gin.Engine {
	wire.Build(
		// db层
		InitDb, InitCache,
		cache.NewUserRedisCache, dao.NewUseMysqlDAO,
		// repository
		repository.NewCachedUserRepository,
		// service
		service.NewUserService,
		// web
		web.NewUserHandler,
		web.InitWebService, web.InitUserMidleware,
	)
	return new(gin.Engine)
}

func InitDownCacheWebService() *gin.Engine {
	wire.Build(
		// db层
		InitDb, InitMemoryCache,
		cache.NewUserMemoryCache, dao.NewUseMysqlDAO,
		// repository
		repository.NewCachedUserRepository,
		// service
		service.NewUserService,
		// web
		web.NewUserHandler,
		web.InitWebService, web.InitUserMidleware,
	)
	return new(gin.Engine)
}
