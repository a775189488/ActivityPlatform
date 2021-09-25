package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ActivityType struct {
	Id         uint64 `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	Parent     uint64 `json:"parent"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
}

func (ActivityType) TableName() string {
	return "act_type_tab"
}

func (a *ActivityType) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (a *ActivityType) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}
