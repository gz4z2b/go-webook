//go:build k8s

/*
 * @Author: p_hanxichen
 * @Date: 2023-09-04 20:42:52
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/conf/k8s.go
 * @Description: k8s配置
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */

package conf

var Db = DbConf{
	Host:     "webook-mysql",
	Port:     "11309",
	Password: "gz4z2b",
	User:     "root",
	Db:       "webook",
}

var Redis = RedisConf{
	Host:     "webook-redis",
	Port:     "11310",
	Password: "",
	Db:       1,
}

var Keys = KeyConf{
	AuthorizationKey: "MXE4iuIoCMBX3Qnco2eqCkSVpIh1v8L3GirpwushYuuhoZI9DoFg7MlJbIYEZmKr",
	EncryptKey:       "he4GdM1Ki9OVbCAqgGCJQeoCffADbx3C",
}
