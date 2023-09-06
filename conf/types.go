/*
 * @Author: p_hanxichen
 * @Date: 2023-09-04 20:24:41
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/conf/types.go
 * @Description:配置
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package conf

type DbConf struct {
	Host     string
	User     string
	Port     string
	Password string
	Db       string
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
}

type KeyConf struct {
	AuthorizationKey string
	EncryptKey       string
}
