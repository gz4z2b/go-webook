/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 17:17:35
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/dao/interface.go
 * @Description: 用户数据库接口
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package dao

import "context"

type UserDAO interface {
	Insert(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id uint64) (User, error)
	FindProfileByUser(ctx context.Context, user User) (Profile, error)
	InsertProfile(ctx context.Context, user User, profile Profile) (Profile, error)
	UpdateProfile(ctx context.Context, profile Profile) (Profile, error)
}
