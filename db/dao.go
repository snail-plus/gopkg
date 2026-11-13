package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// BaseRepo DB通用逻辑接口
type BaseRepo[T schema.Tabler] interface {
	Get(ctx context.Context, id int) (*T, error)                                // 查询
	GetOne(ctx context.Context, condition any, args ...interface{}) (*T, error) // 查询
	Count(ctx context.Context, value interface{}, condition any, args ...interface{}) (int64, error)
	List(ctx context.Context, condition any, page, pageSize int, order any, args ...interface{}) (int64, []*T, error) // 列表
	QueryList(ctx context.Context, condition any, order any, args ...interface{}) ([]*T, error)                       // 列表
	QueryWithScopes(ctx context.Context, dest interface{}, funcs ...func(*gorm.DB) *gorm.DB) error                    // 查询
	Insert(ctx context.Context, value *T) error                                                                       // 新增
	Update(ctx context.Context, value *T) error                                                                       // 修改
	Delete(ctx context.Context, value *T) error
	Save(ctx context.Context, value *T) error
	SaveBatch(ctx context.Context, values []*T) error                                         // 保存数据
	CustomUpdate(ctx context.Context, sql string, args ...interface{}) error                  // 执行sql
	CustomQuery(ctx context.Context, sql string, dest interface{}, args ...interface{}) error // 执行查询sql
	CustomCount(ctx context.Context, sql string, args ...interface{}) (int64, error)

	ExecInTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// 定义context key用于存储事务DB
type txContextKey struct{}

// WithTx 将事务DB存储到context中
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// 从context中获取事务DB
func getTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// BaseRepoImpl DB通用操作结构体
type BaseRepoImpl[T schema.Tabler] struct {
	db *gorm.DB
}

// getDB 获取数据库连接，如果context中有事务DB则优先使用
func (h *BaseRepoImpl[T]) getDB(ctx context.Context) *gorm.DB {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx
	}
	return h.db
}

// NewBaseRepo DB通用操作结构体初始化
func NewBaseRepo[T schema.Tabler](db *gorm.DB) BaseRepo[T] {
	return &BaseRepoImpl[T]{db: db}
}

// Get 根据ID查询
func (h *BaseRepoImpl[T]) Get(ctx context.Context, id int) (*T, error) {
	var result T
	err := h.getDB(ctx).WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *BaseRepoImpl[T]) QueryList(ctx context.Context, condition any, order any, args ...interface{}) ([]*T, error) {
	var data []*T
	var err error

	db := h.getDB(ctx).WithContext(ctx).Model(new(T))
	if condition != nil {
		db = db.Where(condition, args...)
	}

	if order != nil {
		db = db.Order(order)
	}

	err = db.Find(&data).Error
	return data, err
}

func (h *BaseRepoImpl[T]) GetOne(ctx context.Context, condition any, args ...interface{}) (*T, error) {
	var result T
	err := h.getDB(ctx).WithContext(ctx).Model(&result).Where(condition, args...).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List 列表
func (h *BaseRepoImpl[T]) List(ctx context.Context, condition any, page, pageSize int, order any, args ...interface{}) (int64, []*T, error) {
	var count int64
	var err error
	var data []*T

	if order == nil {
		order = "id desc"
	}

	if condition == nil {
		err = h.getDB(ctx).WithContext(ctx).Model(new(T)).Count(&count).Error
		if err != nil {
			return 0, nil, err
		}

		err = h.getDB(ctx).WithContext(ctx).Limit(pageSize).Offset((page - 1) * pageSize).Order(order).Find(&data).Error
		return count, data, err
	} else {
		err = h.getDB(ctx).WithContext(ctx).Model(new(T)).Where(condition, args...).Count(&count).Error
		if err != nil {
			return 0, nil, err
		}

		err = h.getDB(ctx).WithContext(ctx).Limit(pageSize).Offset((page-1)*pageSize).Where(condition, args...).Order(order).Find(&data).Error
		return count, data, err
	}

}

// Insert 新增
func (h *BaseRepoImpl[T]) Insert(ctx context.Context, value *T) error {
	return h.getDB(ctx).WithContext(ctx).Create(value).Error
}

// Update 修改
func (h *BaseRepoImpl[T]) Update(ctx context.Context, value *T) error {
	return h.getDB(ctx).WithContext(ctx).Updates(value).Error
}

func (h *BaseRepoImpl[T]) Save(ctx context.Context, value *T) error {
	return h.getDB(ctx).WithContext(ctx).Save(value).Error
}

// SaveBatch 保存数据
func (h *BaseRepoImpl[T]) SaveBatch(ctx context.Context, values []*T) error {
	return h.getDB(ctx).WithContext(ctx).CreateInBatches(values, 100).Error
}

// Delete 删除
func (h *BaseRepoImpl[T]) Delete(ctx context.Context, value *T) error {
	return h.getDB(ctx).WithContext(ctx).Delete(value).Error
}

// CustomUpdate 执行sql
func (h *BaseRepoImpl[T]) CustomUpdate(ctx context.Context, sql string, args ...interface{}) error {
	return h.getDB(ctx).WithContext(ctx).Exec(sql, args...).Error
}

// CustomQuery 执行查询sql
func (h *BaseRepoImpl[T]) CustomQuery(ctx context.Context, sql string, dest interface{}, args ...interface{}) error {
	return h.getDB(ctx).WithContext(ctx).Raw(sql, args...).Scan(dest).Error
}

func (h *BaseRepoImpl[T]) CustomCount(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var data int64
	if err := h.getDB(ctx).WithContext(ctx).Raw(sql, args...).Count(&data).Error; err != nil {
		return 0, err
	}

	return data, nil
}

func (h *BaseRepoImpl[T]) QueryWithScopes(ctx context.Context, dest interface{}, funcs ...func(*gorm.DB) *gorm.DB) error {
	return h.getDB(ctx).WithContext(ctx).Scopes(funcs...).Find(dest).Error
}

func (h *BaseRepoImpl[T]) Count(ctx context.Context, value interface{}, condition any, args ...interface{}) (int64, error) {
	var data int64
	if err := h.getDB(ctx).WithContext(ctx).Model(value).Where(condition, args...).Count(&data).Error; err != nil {
		return 0, err
	}

	return data, nil
}

// ExecInTx 在事务中执行函数，自动处理事务的提交和回滚（panic 时也会回滚）。
// 事务内的操作请直接使用 fn 回调中的 tx；若需要在 repo 其他方法内复用同一事务，
// 可派生 ctx 后传入 repo 方法，例如 repo.Get(db.WithTx(ctx, tx), id)。
func (h *BaseRepoImpl[T]) ExecInTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return h.getDB(ctx).WithContext(ctx).Transaction(fn)
}

func (h *BaseRepoImpl[T]) GetDb() *gorm.DB {
	return h.db
}
