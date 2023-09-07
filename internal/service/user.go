/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:18:21
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/service/user.go
 * @Description: 用户服务层
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package service

import (
	"context"
	"errors"

	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailConflict   = repository.ErrEmailConflict
	ErrUserNotFound    = repository.ErrUserNotFound
	ErrPasswordInvalid = errors.New("密码不正确")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

/**
 * @description: 注册
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @return {error}
 */
func (svc *UserService) SignUp(ctx context.Context, user *domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

/**
 * @description: 登录
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @return {*domain.User, error}
 */
func (svc *UserService) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	findUser, err := svc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return &domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password))
	if err != nil {
		return &domain.User{}, ErrPasswordInvalid
	}
	return findUser, nil
}

/**
 * @description: 根据email获取用户
 * @param {context.Context} ctx
 * @param {string} email
 * @return {*domain.User, error}
 */
func (svc *UserService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return svc.repo.FindByEmail(ctx, email)
}

/**
 * @description: 根据id获取用户
 * @param {context.Context} ctx
 * @param {uint64} id
 * @return {*domain.User, error}
 */
func (svc *UserService) FindById(ctx context.Context, id uint64) (*domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

/**
 * @description: 根据用户获取档案
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @return {*domain.Profile, error}
 */
func (svc *UserService) FindProfileByUser(ctx context.Context, user *domain.User) (*domain.Profile, error) {
	findProfile, err := svc.repo.FindProfileByUser(ctx, dao.User{
		Id:    user.Id,
		Email: user.Email,
	})
	if err != nil {
		return &domain.Profile{}, err
	}
	return findProfile, nil
}

/**
 * @description: 添加档案
 * @param {context.Context} ctx
 * @param {*domain.User} user
 * @param {*domain.Profile} profile
 * @return {*domain.Profile, error}
 */
func (svc *UserService) AddProfile(ctx context.Context, user *domain.User, profile *domain.Profile) (*domain.Profile, error) {
	return svc.repo.AddProfile(ctx, user, profile)

}
