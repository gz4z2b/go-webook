//go:build !k8s

/*
 * @Author: p_hanxichen
 * @Date: 2023-09-04 20:15:53
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/conf/dev.go
 * @Description: 测试环境配置
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */

package conf

var Db = DbConf{
	Host:     "127.0.0.1",
	Port:     "13316",
	Password: "gz4z2b",
	User:     "root",
	Db:       "webook",
}

var Redis = RedisConf{
	Host:     "127.0.0.1",
	Port:     "13317",
	Password: "",
}

var Keys = KeyConf{
	AuthorizationKey: "MXE4iuIoCMBX3Qnco2eqCkSVpIh1v8L3GirpwushYuuhoZI9DoFg7MlJbIYEZmKr",
	EncryptKey:       "he4GdM1Ki9OVbCAqgGCJQeoCffADbx3C",
}
