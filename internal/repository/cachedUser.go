/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:29:50
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cachedUser.go
 * @Description: 用户数据抽象层
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package repository

import (
	"context"

	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository/cache"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
)

var (
	ErrEmailConflict   = dao.ErrEmailConflict
	ErrUserNotFound    = dao.ErrUserNotFound
	ErrProfileConflict = dao.ErrProfileConflict
	ErrCacheNotExist   = cache.ErrCacheNotExist
)

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}

}

/**
 * @description: 创建用户
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @return {error}
 */
func (r *CachedUserRepository) Create(ctx context.Context, user *domain.User) error {
	userDao, err := r.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
	if err == nil {
		return r.cache.SetUser(ctx, userDao)
	}
	return err
}

/**
 * @description: 根据email获取用户
 * @param {context.Context} ctx
 * @param {string} email
 * @return {*domain.User, error}
 */
func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.cache.FindUserByEmail(ctx, email)
	if err != nil {
		if err == ErrCacheNotExist {
			user, err = r.dao.FindByEmail(ctx, email)
			if err != nil {
				return &domain.User{}, err
			}
			err = r.cache.SetUser(ctx, user)
			if err != nil {
				// TODO 监控
			}
		} else {
			return &domain.User{}, err
		}
	}
	return &domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, err
}

/**
 * @description: 根据id获取用户
 * @param {context.Context} ctx
 * @param {uint64} id
 * @return {*domain.User, error}
 */
func (r *CachedUserRepository) FindById(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := r.cache.FindUserById(ctx, id)
	if err != nil {
		if err == ErrCacheNotExist {
			user, err = r.dao.FindById(ctx, id)
			if err != nil {
				return &domain.User{}, err
			}
			err = r.cache.SetUser(ctx, user)
			if err != nil {
				// TODO 监控
			}
		} else {
			return &domain.User{}, err
		}
	}
	return &domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, err
}

/**
 * @description: 根据用户获取档案
 * @param {context.Context} ctx
 * @param {dao.User} user
 * @return {*domain.Profile, error}
 */
func (r *CachedUserRepository) FindProfileByUser(ctx context.Context, user dao.User) (*domain.Profile, error) {
	profile, err := r.cache.FindProfileByUser(ctx, user)
	if err != nil {
		if err == ErrCacheNotExist {
			profile, err = r.dao.FindProfileByUser(ctx, user)
			if err != nil {
				return &domain.Profile{}, err
			}
			err = r.cache.SetProfile(ctx, profile)
			if err != nil {
				//TODO 监控
			}
		} else {
			return &domain.Profile{}, err
		}
	}
	return &domain.Profile{
		UserId:      profile.UserId,
		NickName:    profile.Nickname,
		BirthDay:    profile.Birthday,
		Description: profile.Description,
	}, err
}

/**
 * @description: 添加用户档案
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @param {*domain.Profile} profile
 * @return {*domain.Profile, error}
 */
func (r *CachedUserRepository) AddProfile(ctx context.Context, user *domain.User, profile *domain.Profile) (*domain.Profile, error) {
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
	profileDao, err := r.dao.InsertProfile(ctx, userDao, profileDao)
	if err != nil {
		if err == ErrProfileConflict {
			profileDao, err = r.dao.FindProfileByUser(ctx, userDao)
			if err != nil {
				return &domain.Profile{}, err
			}
			profileDao, err = r.dao.UpdateProfile(ctx, profileDao)
			if err != nil {
				return &domain.Profile{}, err
			}
		} else {
			return &domain.Profile{}, err
		}
	}
	profile.UserId = user.Id
	user.Profile = *profile
	err = r.cache.SetProfile(ctx, profileDao)
	if err != nil {
		return &domain.Profile{}, err
	}
	return profile, err
}
