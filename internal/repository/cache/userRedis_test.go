/*
 * @Author: p_hanxichen
 * @Date: 2023-09-07 10:01:36
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/cache/userRedis_test.go
 * @Description: 缓存取用户
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	redismocks "github.com/gz4z2b/go-webook/internal/repository/cache/mocks/redismocks"
	"github.com/gz4z2b/go-webook/internal/repository/dao"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

var user = dao.User{
	Id:         1,
	Email:      "gz4z2b@163.com",
	Password:   "19890821Xi_",
	Createtime: time.Now().Unix(),
	Updatetime: time.Now().Unix(),
	Deletetime: 0,
	Profile:    dao.Profile{},
}
var profile = dao.Profile{
	Id:          1,
	UserId:      1,
	Nickname:    "test",
	Birthday:    time.Now().Unix(),
	Description: "test",
	Createtime:  time.Now().Unix(),
	Updatetime:  time.Now().Unix(),
	Deletetime:  0,
}

func TestUserRedisCache_FindUserById(t *testing.T) {
	tests := []struct {
		name     string
		inputId  uint64
		mock     func(ctrl *gomock.Controller) redis.Cmdable
		wantUser dao.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name:    "正常",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)
				result, _ := json.Marshal(dao.User{
					Id:    1,
					Email: "gz4z2b@163.com",
				})
				resultStr := redis.NewStringCmd(context.Background())
				resultStr.SetVal(string(result))

				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultStr)
				return mock
			},
			wantUser: dao.User{
				Id:    1,
				Email: "gz4z2b@163.com",
			},
			wantErr: nil,
		},
		{
			name:    "缓存不存在",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				resultCmd := redis.NewStringCmd(context.Background())
				resultCmd.SetErr(redis.Nil)
				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)
				return mock
			},
			wantUser: dao.User{},
			wantErr:  ErrCacheNotExist,
		},
		{
			name:    "缓存炸了",
			inputId: uint64(1),
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)
				result, _ := json.Marshal(dao.User{})
				resultCmd := redis.NewStringCmd(context.Background())
				resultCmd.SetVal(string(result))
				resultCmd.SetErr(errors.New("缓存炸了"))

				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)
				return mock
			},
			wantUser: dao.User{},
			wantErr:  errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			redis := tt.mock(ctrl)
			cache := NewUserRedisCache(redis)
			user, err := cache.FindUserById(context.Background(), tt.inputId)

			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserRedisCache_FindUserByEmail(t *testing.T) {

	tests := []struct {
		name       string
		inputEmail string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		wantUser   dao.User
		wantErr    error
	}{
		// TODO: Add test cases.
		{
			name:       "正常",
			inputEmail: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cacheMock := redismocks.NewMockCmdable(ctrl)
				result, _ := json.Marshal(user)
				resultCmd := redis.NewStringCmd(context.Background())
				resultCmd.SetVal(string(result))

				cacheMock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return cacheMock
			},
			wantUser: user,
			wantErr:  nil,
		},
		{
			name:       "缓存不存在",
			inputEmail: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cacheMock := redismocks.NewMockCmdable(ctrl)
				//result, _ := json.Marshal(user)
				resultCmd := redis.NewStringCmd(context.Background())
				//resultCmd.SetVal(string(result))
				resultCmd.SetErr(redis.Nil)

				cacheMock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return cacheMock
			},
			wantUser: dao.User{},
			wantErr:  ErrCacheNotExist,
		},
		{
			name:       "缓存炸了",
			inputEmail: "gz4z2b@163.com",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cacheMock := redismocks.NewMockCmdable(ctrl)
				//result, _ := json.Marshal(user)
				resultCmd := redis.NewStringCmd(context.Background())
				//resultCmd.SetVal(string(result))
				resultCmd.SetErr(errors.New("缓存炸了"))

				cacheMock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return cacheMock
			},
			wantUser: dao.User{},
			wantErr:  errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := tt.mock(ctrl)
			cacheMock := NewUserRedisCache(mock)

			user, err := cacheMock.FindUserByEmail(context.Background(), tt.inputEmail)

			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserRedisCache_FindProfileByUser(t *testing.T) {

	type args struct {
		ctx  context.Context
		user dao.User
	}
	tests := []struct {
		name        string
		args        args
		mock        func(ctrl *gomock.Controller) redis.Cmdable
		wantProfile dao.Profile
		wantErr     error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)
				resultCmd := redis.NewStringCmd(context.Background())
				result, _ := json.Marshal(profile)
				resultCmd.SetVal(string(result))
				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return mock
			},
			wantProfile: profile,
			wantErr:     nil,
		},
		{
			name: "缓存不存在",
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)
				resultCmd := redis.NewStringCmd(context.Background())
				//result, _ := json.Marshal(profile)
				//resultCmd.SetVal(string(result))
				resultCmd.SetErr(redis.Nil)
				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return mock
			},
			wantProfile: dao.Profile{},
			wantErr:     ErrCacheNotExist,
		},
		{
			name: "缓存炸了",
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)
				resultCmd := redis.NewStringCmd(context.Background())
				//result, _ := json.Marshal(profile)
				//resultCmd.SetVal(string(result))
				resultCmd.SetErr(errors.New("缓存炸了"))
				mock.EXPECT().Get(context.Background(), gomock.Any()).Return(resultCmd)

				return mock
			},
			wantProfile: dao.Profile{},
			wantErr:     errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := NewUserRedisCache(tt.mock(ctrl))
			got, err := u.FindProfileByUser(tt.args.ctx, tt.args.user)
			assert.Equal(t, tt.wantProfile, got)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}

func TestUserRedisCache_SetUser(t *testing.T) {

	type args struct {
		ctx  context.Context
		user dao.User
	}
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				userStr, _ := json.Marshal(user)
				statusCmd := redis.NewStatusCmd(context.Background())

				mock.EXPECT().Set(context.Background(), gomock.Any(), userStr, gomock.Any()).Return(statusCmd)
				mock.EXPECT().Set(context.Background(), gomock.Any(), userStr, gomock.Any()).Return(statusCmd)

				return mock
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: nil,
		},
		{
			name: "用户缓存炸了",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				userStr, _ := json.Marshal(user)
				statusCmd := redis.NewStatusCmd(context.Background())
				statusCmd.SetErr(errors.New("用户缓存炸了"))

				mock.EXPECT().Set(context.Background(), gomock.Any(), userStr, gomock.Any()).Return(statusCmd)

				return mock
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: errors.New("用户缓存炸了"),
		},
		{
			name: "用户byemail缓存炸了",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				userStr, _ := json.Marshal(user)
				statusCmd := redis.NewStatusCmd(context.Background())
				byEmailStatusCmd := redis.NewStatusCmd(context.Background())
				byEmailStatusCmd.SetErr(errors.New("用户byemail缓存炸了"))

				mock.EXPECT().Set(context.Background(), gomock.Any(), userStr, gomock.Any()).Return(statusCmd)
				mock.EXPECT().Set(context.Background(), gomock.Any(), userStr, gomock.Any()).Return(byEmailStatusCmd)

				return mock
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: errors.New("用户byemail缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := NewUserRedisCache(tt.mock(ctrl))

			err := u.SetUser(tt.args.ctx, tt.args.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserRedisCache_SetProfile(t *testing.T) {
	type args struct {
		ctx     context.Context
		profile dao.Profile
	}
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "正常",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				resultStr, _ := json.Marshal(profile)

				statusCmd := redis.NewStatusCmd(context.Background())
				mock.EXPECT().Set(context.Background(), gomock.Any(), resultStr, gomock.Any()).Return(statusCmd)

				return mock
			},
			args: args{
				ctx:     context.Background(),
				profile: profile,
			},
			wantErr: nil,
		},
		{
			name: "缓存炸了",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mock := redismocks.NewMockCmdable(ctrl)

				resultStr, _ := json.Marshal(profile)

				statusCmd := redis.NewStatusCmd(context.Background())
				statusCmd.SetErr(errors.New("缓存炸了"))
				mock.EXPECT().Set(context.Background(), gomock.Any(), resultStr, gomock.Any()).Return(statusCmd)

				return mock
			},
			args: args{
				ctx:     context.Background(),
				profile: profile,
			},
			wantErr: errors.New("缓存炸了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := NewUserRedisCache(tt.mock(ctrl))
			err := u.SetProfile(tt.args.ctx, tt.args.profile)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
