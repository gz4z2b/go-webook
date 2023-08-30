/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 20:45:58
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/web/middleware/login.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gz4z2b/go-webook/internal/domain"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (loginMiddlewareBuilder *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	loginMiddlewareBuilder.paths = append(loginMiddlewareBuilder.paths, path)
	return loginMiddlewareBuilder
}

func (loginMiddlewareBuilder *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, value := range loginMiddlewareBuilder.paths {
			if value == ctx.Request.URL.Path {
				return
			}
		}

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &domain.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("1J4HLQesjfta8xLQwFDT079VZ6fAasTeyHvlvEMRe4JPVu2DSXJV1OeWflzWJKrv"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || len(claims.Email) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ctx.Request.UserAgent() != claims.UserAgent {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt.Sub(time.Now()) < time.Minute*30 {
			userClaims := domain.UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
				Email:     claims.Email,
				UserAgent: ctx.Request.UserAgent(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
			tokenStr, err := token.SignedString([]byte("1J4HLQesjfta8xLQwFDT079VZ6fAasTeyHvlvEMRe4JPVu2DSXJV1OeWflzWJKrv"))
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("user_email", claims.Email)

		// session := sessions.Default(ctx)
		// session.Options(sessions.Options{
		// 	MaxAge: 30,
		// })
		// email := session.Get("user_email")
		// if email == nil {
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }
		// ctx.Set("user_email", email)
		// now := time.Now().UnixMilli()
		// lastLoginTime := session.Get("login_time")
		// if lastLoginTime == nil {
		// 	session.Set("login_time", now)
		// 	if err := session.Save(); err != nil {
		// 		panic(err)
		// 	}
		// 	return
		// }
		// lastLoginTimeVal, ok := lastLoginTime.(int64)
		// if !ok {
		// 	ctx.AbortWithStatus(http.StatusInternalServerError)
		// }
		// if now-lastLoginTimeVal > 10000 {
		// 	session.Set("login_time", now)
		// 	session.Save()
		// }

	}
}
