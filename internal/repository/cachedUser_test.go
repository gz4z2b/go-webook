/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:29:50
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cachedUser_test.go
 * @Description: 用户数据抽象层
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/repository/cache"
	cachemocks "github.com/gz4z2b/go-webook/internal/repository/cache/mocks"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	daomocks "github.com/gz4z2b/go-webook/internal/repository/dao/mocks"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_Create(t *testing.T) {

	tests := []struct {
		name      string
		inputUser *domain.User
		mock      func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		wantErr   error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			inputUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				cache := cachemocks.NewMockUserCache(ctrl)

				daoMock.EXPECT().Insert(gomock.Any(), dao.User{
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}).Return(dao.User{
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}, nil)
				cache.EXPECT().SetUser(gomock.Any(), dao.User{
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}).Return(nil)

				return daoMock, cache
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			inputUser: &domain.User{
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				//cache := cachemocks.NewMockUserCache(ctrl)

				daoMock.EXPECT().Insert(gomock.Any(), dao.User{
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}).Return(dao.User{}, ErrEmailConflict)
				// cache.EXPECT().SetUser(gomock.Any(), dao.User{
				// 	Email:    "gz4z2b@163.com",
				// 	Password: "19890821Xi_",
				// }).Return(nil)

				return daoMock, nil
			},
			wantErr: ErrEmailConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tt.mock(ctrl)
			repo := NewCachedUserRepository(dao, cache)

			err := repo.Create(context.Background(), tt.inputUser)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCachedUserRepository_FindByEmail(t *testing.T) {

	tests := []struct {
		name       string
		inputEmail string
		mock       func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		wantUser   *domain.User
		wantErr    error
	}{
		// TODO: Add test cases.
		{
			name:       "正常",
			inputEmail: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindByEmail(gomock.Any(), "gz4z2b@163.com").Return(dao.User{
					Id:       1,
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}, nil)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserByEmail(gomock.Any(), "gz4z2b@163.com").Return(dao.User{}, ErrCacheNotExist)
				cacheMock.EXPECT().SetUser(gomock.Any(), dao.User{
					Id:       1,
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}).Return(nil)

				return daoMock, cacheMock
			},
			wantUser: &domain.User{
				Id:       1,
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			wantErr: nil,
		},
		{
			name:       "不存在",
			inputEmail: "gz4z2a@163.com",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindByEmail(gomock.Any(), "gz4z2a@163.com").Return(dao.User{}, ErrUserNotFound)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserByEmail(gomock.Any(), "gz4z2a@163.com").Return(dao.User{}, ErrCacheNotExist)

				return daoMock, cacheMock
			},
			wantUser: &domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name:       "缓存找到",
			inputEmail: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserByEmail(gomock.Any(), "gz4z2b@163.com").Return(dao.User{
					Id:       1,
					Email:    "gz4z2b@163.com",
					Password: "19890821Xi_",
				}, nil)

				return daoMock, cacheMock
			},
			wantUser: &domain.User{
				Id:       1,
				Email:    "gz4z2b@163.com",
				Password: "19890821Xi_",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tt.mock(ctrl)
			repo := NewCachedUserRepository(dao, cache)

			user, err := repo.FindByEmail(context.Background(), tt.inputEmail)

			assert.Equal(t, user, tt.wantUser)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestCachedUserRepository_FindById(t *testing.T) {
	tests := []struct {
		name     string
		inputId  uint64
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		wantUser *domain.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name:    "正常",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindById(gomock.Any(), uint64(1)).Return(dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}, nil)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserById(gomock.Any(), uint64(1)).Return(dao.User{}, ErrCacheNotExist)
				cacheMock.EXPECT().SetUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(nil)

				return daoMock, cacheMock
			},
			wantUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			wantErr: nil,
		},
		{
			name:    "不存在",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindById(gomock.Any(), uint64(1)).Return(dao.User{}, ErrUserNotFound)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserById(gomock.Any(), uint64(1)).Return(dao.User{}, ErrCacheNotExist)

				return daoMock, cacheMock
			},
			wantUser: &domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name:    "缓存获取失败",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserById(gomock.Any(), uint64(1)).Return(dao.User{}, errors.New("缓d炸了"))

				return daoMock, cacheMock
			},
			wantUser: &domain.User{},
			wantErr:  errors.New("缓d炸了"),
		},
		{
			name:    "设置缓存失败",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindById(gomock.Any(), uint64(1)).Return(dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}, nil)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindUserById(gomock.Any(), uint64(1)).Return(dao.User{}, ErrCacheNotExist)
				cacheMock.EXPECT().SetUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(errors.New("设置失败"))

				return daoMock, cacheMock
			},
			wantUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			wantErr: errors.New("设置失败"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tt.mock(ctrl)
			repo := NewCachedUserRepository(dao, cache)

			user, err := repo.FindById(context.Background(), tt.inputId)

			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCachedUserRepository_FindProfileByUser(t *testing.T) {
	tests := []struct {
		name        string
		inputUser   dao.User
		mock        func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		wantProfile *domain.Profile
		wantErr     error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			inputUser: dao.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindProfileByUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(dao.Profile{}, ErrCacheNotExist)
				cacheMock.EXPECT().SetProfile(gomock.Any(), dao.Profile{
					UserId:   1,
					Nickname: "test",
				}).Return(nil)

				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindProfileByUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(dao.Profile{
					UserId:   1,
					Nickname: "test",
				}, nil)

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{
				UserId:   1,
				NickName: "test",
			},
			wantErr: nil,
		},
		{
			name: "不存在",
			inputUser: dao.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindProfileByUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(dao.Profile{}, ErrCacheNotExist)

				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().FindProfileByUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(dao.Profile{}, errors.New("档案不存在"))

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("档案不存在"),
		},
		{
			name: "缓存炸了",
			inputUser: dao.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().FindProfileByUser(gomock.Any(), dao.User{
					Id:    uint64(1),
					Email: "gz4z2b@163.com",
				}).Return(dao.Profile{}, errors.New("缓存炸了"))

				daoMock := daomocks.NewMockUserDAO(ctrl)
				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tt.mock(ctrl)
			repo := NewCachedUserRepository(dao, cache)

			profile, err := repo.FindProfileByUser(context.Background(), tt.inputUser)

			assert.Equal(t, tt.wantProfile, profile)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCachedUserRepository_AddProfile(t *testing.T) {
	tests := []struct {
		name         string
		inputUser    *domain.User
		inputProfile *domain.Profile
		mock         func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		wantProfile  *domain.Profile
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "正常添加",
			inputUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().InsertProfile(context.Background(), gomock.Any(), gomock.Any()).Return(dao.Profile{
					UserId:   uint64(1),
					Nickname: "test",
				}, nil)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().SetProfile(gomock.Any(), gomock.Any()).Return(nil)

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{
				UserId:   uint64(1),
				NickName: "test",
			},
			wantErr: nil,
		},
		{
			name: "更新获取失败",
			inputUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().InsertProfile(context.Background(), gomock.Any(), gomock.Any()).Return(dao.Profile{}, ErrProfileConflict)
				daoMock.EXPECT().FindProfileByUser(gomock.Any(), gomock.Any()).Return(dao.Profile{
					UserId:   uint64(1),
					Nickname: "test",
				}, errors.New("档案不存在"))

				cacheMock := cachemocks.NewMockUserCache(ctrl)

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("档案不存在"),
		},
		{
			name: "更新失败",
			inputUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().InsertProfile(context.Background(), gomock.Any(), gomock.Any()).Return(dao.Profile{}, ErrProfileConflict)
				daoMock.EXPECT().FindProfileByUser(gomock.Any(), gomock.Any()).Return(dao.Profile{
					UserId:   uint64(1),
					Nickname: "test",
				}, nil)
				daoMock.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(dao.Profile{}, errors.New("更新炸了"))

				cacheMock := cachemocks.NewMockUserCache(ctrl)

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("更新炸了"),
		},
		{
			name: "添加炸了",
			inputUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().InsertProfile(context.Background(), gomock.Any(), gomock.Any()).Return(dao.Profile{
					UserId:   uint64(1),
					Nickname: "test",
				}, errors.New("添加炸了"))

				cacheMock := cachemocks.NewMockUserCache(ctrl)

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("添加炸了"),
		},
		{
			name: "缓存炸了",
			inputUser: &domain.User{
				Id:    uint64(1),
				Email: "gz4z2b@163.com",
			},
			inputProfile: &domain.Profile{
				NickName: "test",
			},
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				daoMock.EXPECT().InsertProfile(context.Background(), gomock.Any(), gomock.Any()).Return(dao.Profile{
					UserId:   uint64(1),
					Nickname: "test",
				}, nil)

				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().SetProfile(gomock.Any(), gomock.Any()).Return(errors.New("缓存炸了"))

				return daoMock, cacheMock
			},
			wantProfile: &domain.Profile{},
			wantErr:     errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tt.mock(ctrl)
			repo := NewCachedUserRepository(dao, cache)

			profile, err := repo.AddProfile(context.Background(), tt.inputUser, tt.inputProfile)

			assert.Equal(t, tt.wantProfile, profile)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
