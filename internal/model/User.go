/*
 * @Author: p_hanxichen
 * @Date: 2023-08-22 16:27:35
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/model/User.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       uint   `gorm:"primarykey,autoIncrement`
	Email    string `gorm:"column:email`
	Password string `gorm:"column:password`
}
