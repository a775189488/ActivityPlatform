package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ActivityUser struct {
	Id         uint64 `gorm:"primary_key" json:"id"`
	ActId      uint64 `json:"act_id"`
	UserId     uint64 `json:"user_id"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
}

func (ActivityUser) TableName() string {
	return "act_user_tab"
}

func (u *ActivityUser) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (u *ActivityUser) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}
