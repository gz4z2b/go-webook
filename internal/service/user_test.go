/*
 * @Author: p_hanxichen
 * @Date: 2023-09-20 19:53:16
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/service/user_test.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository"
	repomocks "github.com/gz4z2b/go-webook/internal/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestUserServiceInstance_SignUp(t *testing.T) {
	tests := []struct {
		name      string
		inputUser *domain.User
		mock      func(ctrl *gomock.Controller) repository.UserRepository
		wantErr   error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			inputUser: &domain.User{
				Email:    "gz4z2b@13.com",
				Password: "19890821Xi_",
			},
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return repo
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tt.mock(ctrl))
			err := svc.SignUp(context.Background(), tt.inputUser)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceInstance_Login(t *testing.T) {

	tests := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) repository.UserRepository
		inputUser *domain.User
		wantUser  *domain.User
		wantErr   error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				user := &domain.User{
					Email:    "gz4z2b@163.com",
					Password: "$2a$10$2/zoj94WMfc7xvGv9NNsmuptftGX3MnyBiycLYc0lmYsKrGJGOkNK",
				}
				repo.EXPECT().FindByEmail(gomock.Any(), "gz4z2b@163.com").Return(user, nil)
				return repo
			},
			inputUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			wantUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "$2a$10$2/zoj94WMfc7xvGv9NNsmuptftGX3MnyBiycLYc0lmYsKrGJGOkNK",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "gz4z2b@163.com").Return(&domain.User{}, ErrUserNotFound)
				return repo
			},
			inputUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			wantUser: &domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				user := &domain.User{
					Email:    "gz4z2b@163.com",
					Password: "$2a$10$2/zoj94WMfc7xvGv9NNsmuptftGX3MnyBiycLYc0lmYsKrGJGOkNK",
				}
				repo.EXPECT().FindByEmail(gomock.Any(), "gz4z2b@163.com").Return(user, nil)
				return repo
			},
			inputUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi",
			},
			wantUser: &domain.User{},
			wantErr:  ErrPasswordInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tt.mock(ctrl))
			user, err := svc.Login(context.Background(), tt.inputUser)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUser, user)
		})
	}
}

func TestUserServiceInstance_FindByEmail(t *testing.T) {

	tests := []struct {
		name     string
		email    string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		wantUser *domain.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name:  "正常",
			email: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "gz4z2b@163.com").Return(&domain.User{
					Email: "gz4z2b@163.com",
				}, nil)
				return repo
			},
			wantUser: &domain.User{
				Email: "gz4z2b@163.com",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tt.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.FindByEmail(context.Background(), tt.email)

			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceInstance_FindById(t *testing.T) {

	tests := []struct {
		name     string
		inputId  uint64
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		wantUser *domain.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name:    "正常",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(&domain.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}, nil)
				return repo
			},
			wantUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tt.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.FindById(context.Background(), tt.inputId)

			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceInstance_FindProfileByUser(t *testing.T) {
	tests := []struct {
		name        string
		inputUser   *domain.User
		mock        func(ctrl *gomock.Controller) repository.UserRepository
		wantProfile *domain.Profile
		wantErr     error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			inputUser: &domain.User{
				Id:    1,
				Email: "gz4z2b@163.com",
			},
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindProfileByUser(gomock.Any(), gomock.Any()).Return(&domain.Profile{
					UserId:   1,
					NickName: "test",
				}, nil)
				return repo
			},
			wantProfile: &domain.Profile{
				UserId:   1,
				NickName: "test",
			},
			wantErr: nil,
		},
		{
			name: "档案不存在",
			inputUser: &domain.User{
				Id:    1,
				Email: "gz4z2b@163.com",
			},
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindProfileByUser(gomock.Any(), gomock.Any()).Return(&domain.Profile{}, errors.New("档案不存在"))
				return repo
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("档案不存在"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tt.mock(ctrl)
			svc := NewUserService(repo)
			profile, err := svc.FindProfileByUser(context.Background(), tt.inputUser)

			assert.Equal(t, tt.wantProfile, profile)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceInstance_AddProfile(t *testing.T) {
	tests := []struct {
		name         string
		inputUser    *domain.User
		inputProfile *domain.Profile
		mock         func(ctrl *gomock.Controller) repository.UserRepository
		wantProfile  *domain.Profile
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			inputUser: &domain.User{
				Id:    1,
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().AddProfile(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Profile{
					UserId:   1,
					NickName: "test",
				}, nil)
				return repo
			},
			wantProfile: &domain.Profile{
				UserId:   1,
				NickName: "test",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tt.mock(ctrl)
			svc := NewUserService(repo)
			profile, err := svc.AddProfile(context.Background(), tt.inputUser, tt.inputProfile)

			assert.Equal(t, tt.wantProfile, profile)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
