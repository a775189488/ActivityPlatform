package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ActivityComment struct {
	Id         uint64 `gorm:"primary_key" json:"id"`
	ActId      uint64 `json:"act_id"`
	UserId     uint64 `json:"user_id"`
	Message    string `json:"message"`
	Parent     uint64 `json:"parent"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
}

func (ActivityComment) TableName() string {
	return "comment_tab"
}

func (a *ActivityComment) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (a *ActivityComment) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}
