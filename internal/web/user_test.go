/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 17:45:51
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/web/user_test.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gz4z2b/go-webook/internal/domain"
	"github.com/gz4z2b/go-webook/internal/service"
	svcmocks "github.com/gz4z2b/go-webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestEncrypt(t *testing.T) {
	password := "hello#123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestTest(t *testing.T) {
	str := fmt.Sprintf("%s:%s", "gz3z2b", "db")
	println(str)
	assert.NoError(t, nil)
}

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		mock     func(ctrl *gomock.Controller) service.UserService
		wantCode int
		wantBody string
	}{
		{
			name: "正常流程",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_",
				"confirmPassword": "19890821Xi_"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				service.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return service
			},
			wantCode: http.StatusOK,
			wantBody: "success",
		},
		{
			name: "数据格式错误",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_",
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				return service
			},
			wantCode: http.StatusBadRequest,
			wantBody: "数据格式错误",
		},
		{
			name: "两次输入的密码不一致",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_",
				"confirmPassword": "19890821Xi"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				return service
			},
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码复杂度不够",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi",
				"confirmPassword": "19890821Xi"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				return service
			},
			wantCode: http.StatusOK,
			wantBody: "密码复杂度不够",
		},
		{
			name: "邮箱格式不正确",
			input: []byte(`{
				"email": "gz4z2b@163",
				"password": "19890821Xi_",
				"confirmPassword": "19890821Xi_"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				return service
			},
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不正确",
		},
		{
			name: "邮箱冲突啦~~",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_",
				"confirmPassword": "19890821Xi_"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrEmailConflict)
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突啦~~",
		},
		{
			name: "service出错",
			input: []byte(`{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_",
				"confirmPassword": "19890821Xi_"
			}`),
			mock: func(ctrl *gomock.Controller) service.UserService {
				service := svcmocks.NewMockUserService(ctrl)
				service.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("系统出错"))
				return service
			},
			wantCode: http.StatusOK,
			wantBody: "系统出错",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer(tc.input))
			// require.NoError(t, err)
			resp := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewUserHandler(tc.mock(ctrl))
			server := InitWebService(handler, []gin.HandlerFunc{})
			server.ServeHTTP(resp, req)

			assert.Equal(t, resp.Code, tc.wantCode)
			assert.Equal(t, resp.Body.String(), tc.wantBody)

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := svcmocks.NewMockUserService(ctrl)
	service.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("test"))

	err := service.SignUp(context.Background(), &domain.User{})
	t.Log(err)
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mock     func(ctrl *gomock.Controller) service.UserService
		wantCode int
		wantBody string
	}{
		{
			name: "正常",
			input: `{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_"
			}`,
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&domain.User{}, nil)
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: "登录成功",
		},
		{
			name: "输入数据格式错误",
			input: `{
				"email": "gz4z2b@163.com",
			}`,
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			wantCode: http.StatusBadRequest,
			wantBody: "输入数据格式错误",
		},
		{
			name: "密码错误",
			input: `{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi"
			}`,
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&domain.User{}, service.ErrPasswordInvalid)
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: "邮箱或密码错误",
		},
		{
			name: "系统错误",
			input: `{
				"email": "gz4z2b@163.com",
				"password": "19890821Xi_"
			}`,
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&domain.User{}, errors.New("系统错误"))
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer([]byte(tt.input)))
			resp := httptest.NewRecorder()

			handler := NewUserHandler(tt.mock(ctrl))
			server := InitWebService(handler, InitUserMidleware())
			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())

		})
	}
}

func TestUserHandler_Edit(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		mock     func(ctrl *gomock.Controller) service.UserService
		wantCode int
		wantBody string
	}{
		{
			name: "正常",
			input: `{
				"nick_name": "陈瀚禧",
				"birth_day": "1989-08-21",
				"description": "简介"
			}`,
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().AddProfile(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Profile{}, nil)
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: "修改成功",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(http.MethodPost, "/users/edit", bytes.NewBuffer([]byte(tt.input)))
			resp := httptest.NewRecorder()

			handler := NewUserHandler(tt.mock(ctrl))
			server := InitWebService(handler, []gin.HandlerFunc{})
			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantBody, resp.Body.String())
			assert.Equal(t, tt.wantCode, resp.Code)
		})
	}
}
