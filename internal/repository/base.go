package repository

import (
	"entrytask/internal/common/db"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/utils"
	"github.com/jinzhu/gorm"
)

type BaseRepo struct {
	Source db.IDb         `inject:""`
	Log    logger.ILogger `inject:""`
}

// Create 创建实体
func (b *BaseRepo) Create(value interface{}) error {
	return b.Source.DB().Create(value).Error
}

// Save 保存实体
func (b *BaseRepo) Save(value interface{}) error {
	return b.Source.DB().Save(value).Error
}

// Updates 更新实体
func (b *BaseRepo) Updates(model interface{}, value interface{}) error {
	return b.Source.DB().Model(model).Updates(value).Error
}

// DeleteByWhere 根据条件删除实体
func (b *BaseRepo) DeleteByWhere(model, where interface{}) (count int64, err error) {
	db := b.Source.DB().Where(where).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DeleteByID 根据id删除实体
func (b *BaseRepo) DeleteByID(model interface{}, id int) error {
	return b.Source.DB().Where("id=?", id).Delete(model).Error
}

// Find 根据条件返回数据
func (b *BaseRepo) Find(where interface{}, out interface{}, sel string, orders ...string) error {
	db := b.Source.DB().Where(where)
	if sel != "" {
		db = db.Select(sel)
	}
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// GetPages 分页返回数据
func (b *BaseRepo) GetPages(model interface{}, out interface{}, pageIndex int32, pageSize int32, totalCount *int32, where interface{}, orders ...string) error {
	db := b.Source.DB().Model(model).Where(model)
	where_ := make(map[string]interface{})
	if where != nil {
		for k, v := range where.(map[string]interface{}) {
			where_[utils.Camel2Case(k)] = v
		}
	}
	db = db.Where(where_)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	err := db.Count(totalCount).Error
	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

func (b *BaseRepo) FirstByID(out interface{}, id uint64) error {
	return b.Source.DB().First(out, id).Error
}

//GetTransaction 获取事务
func (b *BaseRepo) GetTransaction() *gorm.DB {
	return b.Source.DB().Begin()
}

func (b *BaseRepo) First(where interface{}, out interface{}, selects ...string) error {
	db := b.Source.DB().Where(where)
	if len(selects) > 0 {
		for _, sel := range selects {
			db = db.Select(sel)
		}
	}
	return db.First(out).Error
}
