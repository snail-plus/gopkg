package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// BaseRepo DB通用逻辑接口
type BaseRepo[T schema.Tabler] interface {
	Get(ctx context.Context, id int) (*T, error)                                                           // 查询
	GetOne(ctx context.Context, condition any, args ...interface{}) (*T, error)                            // 查询
	List(ctx context.Context, condition any, page, pageSize int, args ...interface{}) (int64, []*T, error) // 列表
	QueryWithScopes(ctx context.Context, dest interface{}, funcs ...func(*gorm.DB) *gorm.DB) error         // 查询
	Insert(ctx context.Context, value *T) error                                                            // 新增
	Update(ctx context.Context, value *T) error                                                            // 修改
	Delete(ctx context.Context, value *T) error                                                            // 删除
	CustomUpdate(ctx context.Context, sql string, args ...interface{}) error                               // 执行sql
	CustomQuery(ctx context.Context, sql string, dest interface{}, args ...interface{}) error              // 执行查询sql
	CustomCount(ctx context.Context, sql string, args ...interface{}) (int64, error)
}

// BaseRepoImpl DB通用操作结构体
type BaseRepoImpl[T schema.Tabler] struct {
	db *gorm.DB
}

// NewBaseRepo DB通用操作结构体初始化
func NewBaseRepo[T schema.Tabler](db *gorm.DB) BaseRepo[T] {
	return &BaseRepoImpl[T]{db: db}
}

// Get 根据ID查询
func (h *BaseRepoImpl[T]) Get(ctx context.Context, id int) (*T, error) {
	var result T
	err := h.db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *BaseRepoImpl[T]) GetOne(ctx context.Context, condition any, args ...interface{}) (*T, error) {
	var result T
	err := h.db.WithContext(ctx).Model(&result).Where(condition, args...).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List 列表
func (h *BaseRepoImpl[T]) List(ctx context.Context, condition any, page, pageSize int, args ...interface{}) (int64, []*T, error) {
	var count int64
	var tmp *T
	var err error
	var data []*T

	if condition == nil {
		err = h.db.WithContext(ctx).Model(tmp).Count(&count).Error
		if err != nil {
			return 0, nil, err
		}

		err = h.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&data).Error
		return count, data, err
	} else {
		err = h.db.WithContext(ctx).Model(tmp).Where(condition, args...).Count(&count).Error
		if err != nil {
			return 0, nil, err
		}

		err = h.db.Limit(pageSize).Offset((page-1)*pageSize).Where(condition, args...).Find(&data).Error
		return count, data, err
	}

}

// Insert 新增
func (h *BaseRepoImpl[T]) Insert(ctx context.Context, value *T) error {
	return h.db.WithContext(ctx).Create(value).Error
}

// Update 修改
func (h *BaseRepoImpl[T]) Update(ctx context.Context, value *T) error {
	return h.db.WithContext(ctx).Updates(value).Error
}

// Delete 删除
func (h *BaseRepoImpl[T]) Delete(ctx context.Context, value *T) error {
	return h.db.WithContext(ctx).Delete(value).Error
}

// CustomUpdate 执行sql
func (h *BaseRepoImpl[T]) CustomUpdate(ctx context.Context, sql string, args ...interface{}) error {
	return h.db.WithContext(ctx).Exec(sql, args...).Error
}

// CustomQuery 执行sql
func (h *BaseRepoImpl[T]) CustomQuery(ctx context.Context, sql string, dest interface{}, args ...interface{}) error {
	return h.db.WithContext(ctx).Raw(sql, args...).Scan(dest).Error
}

func (h *BaseRepoImpl[T]) CustomCount(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var data int64
	if err := h.db.WithContext(ctx).Raw(sql, args...).Count(&data).Error; err != nil {
		return 0, err
	}

	return data, nil
}

func (h *BaseRepoImpl[T]) QueryWithScopes(ctx context.Context, dest interface{}, funcs ...func(*gorm.DB) *gorm.DB) error {
	return h.db.WithContext(ctx).Scopes(funcs...).Find(&dest).Error
}

func (h *BaseRepoImpl[T]) GetDb() *gorm.DB {
	return h.db
}
