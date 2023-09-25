/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 10:58:25
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/ioc/cache.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package ioc

import (
	"fmt"

	"github.com/coocood/freecache"
	"github.com/gz4z2b/go-webook/conf"
	"github.com/redis/go-redis/v9"
)

func InitCache() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password, // no password set
		DB:       conf.Redis.Db,       // use default DB
	})
}

func InitMemoryCache() *freecache.Cache {
	return freecache.NewCache(10 * 1024 * 1034)
}
