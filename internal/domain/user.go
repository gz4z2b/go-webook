/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:23:43
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/domain/user.go
 * @Description: 用户领域
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package domain

import jwt "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       uint64  `json:"id"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Profile  Profile `json:"profile"`
}

type Profile struct {
	UserId      uint64 `json:"user_id"`
	NickName    string `json:"nick_name"`
	BirthDay    int64  `json:"birth_day"`
	Description string `json:"description"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email     string
	UserAgent string
}
