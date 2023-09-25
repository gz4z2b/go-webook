/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:33:57
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/repository/dao/userMysql.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package dao

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrEmailConflict   error = errors.New("邮箱冲突")
	ErrUserNotFound    error = gorm.ErrRecordNotFound
	ErrProfileNotFound error = gorm.ErrRecordNotFound
	ErrProfileConflict error = errors.New("个人档案冲突")
)

type UserMysqlDAO struct {
	db *gorm.DB
}

func NewUseMysqlDAO(db *gorm.DB) UserDAO {
	return &UserMysqlDAO{
		db: db,
	}
}

/**
 * @description: 插入用户
 * @param {context.Context} ctx
 * @param {User} user
 * @return {error}
 */
func (u *UserMysqlDAO) Insert(ctx context.Context, user User) (User, error) {
	//now := time.Now().UnixMilli()
	//user.Createtime = now
	//user.Updatetime = now

	err := u.db.WithContext(ctx).Create(&user).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrorNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrorNo {
			return User{}, ErrEmailConflict
		}
	}
	return user, err
}

/**
 * @description: 通过email查询用户
 * @param {context.Context} ctx
 * @param {string} email
 * @return {User, error}
 */
func (u *UserMysqlDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return user, err
}

func (u *UserMysqlDAO) FindById(ctx context.Context, id uint64) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return user, err
}

/**
 * @description: 通过用户查找档案
 * @param {context.Context} ctx
 * @param {User} user
 * @return {Profile, error}
 */
func (u *UserMysqlDAO) FindProfileByUser(ctx context.Context, user User) (Profile, error) {
	var profile Profile
	err := u.db.WithContext(ctx).Where("user_id = ?", user.Id).First(&profile).Error
	if err != nil {
		return Profile{}, err
	}
	return profile, err
}

/**
 * @description: 插入档案
 * @param {context.Context} ctx
 * @param {User} user
 * @param {Profile} profile
 * @return {Profile, error}
 */
func (u *UserMysqlDAO) InsertProfile(ctx context.Context, user User, profile Profile) (Profile, error) {
	profile.UserId = user.Id
	err := u.db.WithContext(ctx).Create(&profile).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrorNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrorNo {
			return Profile{}, ErrProfileConflict
		}
		return Profile{}, err
	}
	return profile, err
}

/**
 * @description: 更新档案
 * @param {context.Context} ctx
 * @param {Profile} profile
 * @return {Profile, error}
 */
func (u *UserMysqlDAO) UpdateProfile(ctx context.Context, profile Profile) (Profile, error) {
	err := u.db.WithContext(ctx).Save(&profile).Error
	return profile, err
}

type User struct {
	Id       uint64 `gorm:"primaryKey,not null,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Createtime int64 `gorm:"autoCreateTime:milli"`
	Updatetime int64 `gorm:"autoUpdateTime:milli"`
	Deletetime int64

	Profile Profile `gorm:"foreignKey:UserId"`
}

// TableName 会将 User 的表名重写为 `profiles`
func (u User) TableName() string {
	return "t_user"
}

type Profile struct {
	Id          uint64 `gorm:"primaryKey,not null,autoIncrement"`
	UserId      uint64 `gorm:"unique"`
	Nickname    string
	Birthday    int64
	Description string
	Createtime  int64 `gorm:"autoCreateTime:milli"`
	Updatetime  int64 `gorm:"autoUpdateTime:milli"`
	Deletetime  int64
}

func (p Profile) TableName() string {
	return "t_user_profile"
}
