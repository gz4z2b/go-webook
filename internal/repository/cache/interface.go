/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 17:05:51
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cache/interface.go
 * @Description: 用户缓存接口
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package cache

import (
	"context"
	"errors"

	"github.com/gz4z2b/go-webook/internal/repository/dao"
)

var (
	ErrCacheNotExist error = errors.New("缓存key不存在")
)

type UserCache interface {
	FindUserById(ctx context.Context, id uint64) (dao.User, error)
	FindUserByEmail(ctx context.Context, email string) (dao.User, error)
	FindProfileByUser(ctx context.Context, user dao.User) (dao.Profile, error)
	SetUser(ctx context.Context, user dao.User) error
	SetProfile(ctx context.Context, profile dao.Profile) error
}
