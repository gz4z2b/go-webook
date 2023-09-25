/*
 * @Author: p_hanxichen
 * @Date: 2023-09-07 10:01:36
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cache/userRedis.go
 * @Description: 缓存取用户
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gz4z2b/go-webook/internal/repository/dao"
	redis "github.com/redis/go-redis/v9"
)

type UserRedisCache struct {
	cache      redis.Cmdable
	expiretion time.Duration
}

func NewUserRedisCache(client redis.Cmdable) UserCache {
	return &UserRedisCache{
		cache:      client,
		expiretion: time.Minute * 15,
	}
}

/**
 * @description: 按照id获取用户
 * @param {context.Context} ctx
 * @param {uint64} id
 * @return {dao.User, errror}
 */
func (u *UserRedisCache) FindUserById(ctx context.Context, id uint64) (dao.User, error) {
	key := u.getUserCacheKey(id)
	result, err := u.cache.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return dao.User{}, ErrCacheNotExist
		}
		return dao.User{}, err
	}
	var user dao.User
	err = json.Unmarshal(result, &user)
	return user, err
}

/**
 * @description: 根据email获取用户
 * @param {context.Context} ctx
 * @param {string} email
 * @return {dao.User, error}
 */
func (u *UserRedisCache) FindUserByEmail(ctx context.Context, email string) (dao.User, error) {
	key := u.getUserCacheEmailKey(email)
	result, err := u.cache.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return dao.User{}, ErrCacheNotExist
		}
		return dao.User{}, err
	}
	var user dao.User
	err = json.Unmarshal(result, &user)
	return user, err
}

/**
 * @description: 通过用户查询档案
 * @param {context.Context} ctx
 * @param {dao.User} user
 * @return {dao.Profile}
 */
func (u *UserRedisCache) FindProfileByUser(ctx context.Context, user dao.User) (dao.Profile, error) {
	key := u.getProfileCacheUserKey(user.Id)
	result, err := u.cache.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return dao.Profile{}, ErrCacheNotExist
		}
		return dao.Profile{}, err
	}
	var profile dao.Profile
	err = json.Unmarshal(result, &profile)
	return profile, err
}

/**
 * @description: 设置缓存
 * @param {context.Context} ctx
 * @param {domain.User} user
 * @return {error}
 */
func (u *UserRedisCache) SetUser(ctx context.Context, user dao.User) error {
	setStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = u.cache.Set(ctx, u.getUserCacheKey(user.Id), setStr, u.expiretion).Err()
	if err != nil {
		return err
	}
	err = u.cache.Set(ctx, u.getUserCacheEmailKey(user.Email), setStr, u.expiretion).Err()
	if err != nil {
		return err
	}
	return nil
}

/**
 * @description: 设置个人档案缓存
 * @param {context.Context} ctx
 * @param {domain.Profile} profile
 * @return {error}
 */
func (u *UserRedisCache) SetProfile(ctx context.Context, profile dao.Profile) error {
	setStr, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	err = u.cache.Set(ctx, u.getProfileCacheUserKey(profile.UserId), setStr, u.expiretion).Err()
	if err != nil {
		return err
	}
	return nil
}

/**
 * @description: 用户信息缓存key
 * @param {uint64} id
 * @return {string}
 */
func (u *UserRedisCache) getUserCacheKey(id uint64) string {
	return fmt.Sprintf("webook:user:getusercachekey:%d", id)
}
func (u *UserRedisCache) getUserCacheEmailKey(email string) string {
	return fmt.Sprintf("webook:user:getusercacheemailkey:%s", email)
}
func (u *UserRedisCache) getProfileCacheUserKey(userId uint64) string {
	return fmt.Sprintf("webook:user:getprofilecacheuserkey:%d", userId)
}
