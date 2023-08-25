/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:24:53
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/main.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package main

import "github.com/gz4z2b/go-webook/internal/web"

func main() {

	server := web.RegisterRoutes()

	// db, err := gorm.Open(mysql.Open("root:gz4z2b@tcp(127.0.0.1:3306)/webook"), &gorm.Config{})
	// if err != nil {
	// 	panic("数据库链接失败")
	// }
	// db = db.Debug()

	// db.AutoMigrate(&User{})

	server.Run(":8080")
}
