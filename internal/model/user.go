package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id         uint64 `gorm:"primary_key" json:"id"`
	Aliasname  string `json:"aliasname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Headpic    string `json:"headpic"`
	Role       int    `json:"role"`
	IsOnline   bool   `json:"is_online"`
	IsEnable   bool   `json:"is_enable"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
	DeleteTime uint64 `json:"delete_time"`
}

func (User) TableName() string {
	return "user_tab"
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (u *User) CompareAndSwap(user *User) bool {
	isSwap := false
	if u.Aliasname != user.Aliasname {
		u.Aliasname = user.Aliasname
		isSwap = true
	}
	if u.Email != user.Email {
		u.Email = user.Email
		isSwap = true
	}
	if u.Headpic != user.Headpic {
		u.Headpic = user.Headpic
		isSwap = true
	}
	return isSwap
}
