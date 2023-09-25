/*
 * @Author: p_hanxichen
 * @Date: 2023-09-15 17:23:01
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/interface.go
 * @Description: 用户数据操作接口
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package repository

import (
	"context"

	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindById(ctx context.Context, id uint64) (*domain.User, error)
	FindProfileByUser(ctx context.Context, user dao.User) (*domain.Profile, error)
	AddProfile(ctx context.Context, user *domain.User, profile *domain.Profile) (*domain.Profile, error)
}
