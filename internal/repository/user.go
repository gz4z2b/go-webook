/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:29:50
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/repository/user.go
 * @Description: 用户数据抽象层
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package repository

import (
	"context"

	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
)

var (
	ErrEmailConflict   = dao.ErrEmailConflict
	ErrEmailNotFound   = dao.ErrEmailNotFound
	ErrProfileConflict = dao.ErrProfileConflict
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.dao.FindByEmail(ctx, email)
}

func (r *UserRepository) FindProfileByUser(ctx context.Context, user *domain.User) (*domain.Profile, error) {
	return r.dao.FindProfileDomainByUser(ctx, dao.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	})
}

func (r *UserRepository) AddProfile(ctx context.Context, user *domain.User, profile *domain.Profile) (*domain.Profile, error) {
	userDao := dao.User{
		Id:    user.Id,
		Email: user.Email,
	}
	profileDao := dao.Profile{
		UserId:      user.Id,
		Nickname:    profile.NickName,
		Birthday:    profile.BirthDay,
		Description: profile.Description,
	}
	profileAdd, err := r.dao.InsertProfile(ctx, userDao, profileDao)
	if err != nil {
		if err == ErrProfileConflict {
			profileDao, err := r.dao.FindProfileByUser(ctx, userDao)
			if err != nil {
				return profileAdd, err
			}
			profileDao.Nickname = profile.NickName
			profileDao.Birthday = profile.BirthDay
			profileDao.Description = profile.Description
			return r.dao.UpdateProfile(ctx, profileDao)
		}
		return profileAdd, err
	}
	return profileAdd, err
}
