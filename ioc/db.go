/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 10:56:44
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/ioc/db.go
 * @Description: 数据库初始化
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package ioc

import (
	"fmt"

	"github.com/gz4z2b/go-webook/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Db.User, conf.Db.Password, conf.Db.Host, conf.Db.Port, conf.Db.Db)), &gorm.Config{})
	if err != nil {
		panic("数据库初始化失败")
	}
	return db
}
