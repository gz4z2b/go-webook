/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:18:21
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/service/user.go
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
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailConflict   = repository.ErrEmailConflict
	ErrEmailNotFound   = repository.ErrEmailNotFound
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

func (svc *UserService) SignUp(ctx context.Context, user *domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	findUser, err := svc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return &domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password))
	if err != nil {
		return &domain.User{}, err
	}
	return findUser, nil
}

func (svc *UserService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	findUser, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return &domain.User{}, err
	}
	return findUser, nil
}

func (svc *UserService) FindProfileByUser(ctx context.Context, user *domain.User) (*domain.Profile, error) {
	findProfile, err := svc.repo.FindProfileByUser(ctx, user)
	if err != nil {
		return &domain.Profile{}, err
	}
	return findProfile, nil
}

func (svc *UserService) AddProfile(ctx context.Context, user *domain.User, profile *domain.Profile) (*domain.Profile, error) {
	return svc.repo.AddProfile(ctx, user, profile)

}
