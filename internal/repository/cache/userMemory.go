/*
 * @Author: p_hanxichen
 * @Date: 2023-09-07 10:01:36
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cache/userMemory.go
 * @Description: 本地缓存取用户
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coocood/freecache"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
)

type UserMemoryCache struct {
	cache      *freecache.Cache
	expiretion time.Duration
}

func NewUserMemoryCache(client *freecache.Cache) UserCache {
	return &UserMemoryCache{
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
func (u *UserMemoryCache) FindUserById(ctx context.Context, id uint64) (dao.User, error) {
	key := u.getUserCacheKey(id)
	result, err := u.cache.Get(key)
	if err != nil {
		if err == freecache.ErrNotFound {
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
func (u *UserMemoryCache) FindUserByEmail(ctx context.Context, email string) (dao.User, error) {
	key := u.getUserCacheEmailKey(email)
	result, err := u.cache.Get(key)
	if err != nil {
		if err == freecache.ErrNotFound {
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
func (u *UserMemoryCache) FindProfileByUser(ctx context.Context, user dao.User) (dao.Profile, error) {
	key := u.getProfileCacheUserKey(user.Id)
	result, err := u.cache.Get(key)
	if err != nil {
		if err == freecache.ErrNotFound {
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
func (u *UserMemoryCache) SetUser(ctx context.Context, user dao.User) error {
	setStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = u.cache.Set(u.getUserCacheKey(user.Id), setStr, int(u.expiretion.Seconds()))
	if err != nil {
		return err
	}
	err = u.cache.Set(u.getUserCacheEmailKey(user.Email), setStr, int(u.expiretion.Seconds()))
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
func (u *UserMemoryCache) SetProfile(ctx context.Context, profile dao.Profile) error {
	setStr, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	err = u.cache.Set(u.getProfileCacheUserKey(profile.UserId), setStr, int(u.expiretion.Seconds()))
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
func (u *UserMemoryCache) getUserCacheKey(id uint64) []byte {
	return []byte(fmt.Sprintf("webook:user:getusercachekey:%d", id))
}
func (u *UserMemoryCache) getUserCacheEmailKey(email string) []byte {
	return []byte(fmt.Sprintf("webook:user:getusercacheemailkey:%s", email))
}
func (u *UserMemoryCache) getProfileCacheUserKey(userId uint64) []byte {
	return []byte(fmt.Sprintf("webook:user:getprofilecacheuserkey:%d", userId))
}
