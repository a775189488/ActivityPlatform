package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Activity struct {
	Id          uint64 `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	Description string `json:"description"`
	Creator     uint64 `json:"creator"`
	// todo 需要转换成实际的活动类型
	ActType uint64 `json:"act_type"`
	Address string `json:"address"`
	status  int    `json:"status"`
	// todo 需要将这两个时间转成time返回给前端
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
}

type ActivityDetail struct {
	Activity
	ActTypeName string `json:"act_type_name"`
}

func (Activity) TableName() string {
	return "act_tab"
}

func (a *Activity) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (a *Activity) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}
