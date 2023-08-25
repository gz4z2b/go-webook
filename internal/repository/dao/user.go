/*
 * @Author: p_hanxichen
 * @Date: 2023-08-23 10:33:57
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/internal/repository/dao/user.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package dao

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/gz4z2b/go-webook/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrEmailConflict   error = errors.New("邮箱冲突")
	ErrEmailNotFound   error = gorm.ErrRecordNotFound
	ErrProfileNotFound error = gorm.ErrRecordNotFound
	ErrProfileConflict error = errors.New("个人档案冲突")
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	//now := time.Now().UnixMilli()
	//user.Createtime = now
	//user.Updatetime = now

	err := u.db.WithContext(ctx).Create(&user).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrorNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrorNo {
			return ErrEmailConflict
		}
	}
	return err
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return &domain.User{}, ErrEmailNotFound
	}
	return &domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, err
}

func (u *UserDAO) FindProfileByUser(ctx context.Context, user User) (Profile, error) {
	var profile Profile
	err := u.db.WithContext(ctx).Where("user_id = ?", user.Id).First(&profile).Error
	if err != nil {
		return Profile{}, err
	}
	return profile, err
}

func (u *UserDAO) FindProfileDomainByUser(ctx context.Context, user User) (*domain.Profile, error) {
	profile, err := u.FindProfileByUser(ctx, user)
	if err != nil {
		return &domain.Profile{}, err
	}
	return &domain.Profile{
		UserId:      profile.UserId,
		NickName:    profile.Nickname,
		BirthDay:    profile.Birthday,
		Description: profile.Description,
	}, nil
}

func (u *UserDAO) InsertProfile(ctx context.Context, user User, profile Profile) (*domain.Profile, error) {
	profile.UserId = user.Id
	err := u.db.WithContext(ctx).Create(&profile).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrorNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrorNo {
			return &domain.Profile{}, ErrProfileConflict
		}
		return &domain.Profile{}, err
	}
	return &domain.Profile{
		UserId:      profile.UserId,
		NickName:    profile.Nickname,
		BirthDay:    profile.Birthday,
		Description: profile.Description,
	}, err
}

func (u *UserDAO) UpdateProfile(ctx context.Context, profile Profile) (*domain.Profile, error) {
	err := u.db.WithContext(ctx).Save(profile).Error
	return &domain.Profile{
		UserId:      profile.UserId,
		NickName:    profile.Nickname,
		BirthDay:    profile.Birthday,
		Description: profile.Description,
	}, err
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
