package query

import (
	"gorm.io/gorm"
)

// BaseModel 只持有原始db，充当基础模型，无任何链式条件方法，安全嵌入
type BaseModel[T any] struct {
	db *gorm.DB
}

func NewBaseModel[T any](db *gorm.DB) *BaseModel[T] {
	return &BaseModel[T]{db: db}
}

// NewQuery 只负责新建临时QueryBuilder
func (b *BaseModel[T]) NewQuery() *QueryBuilder[T] {
	return NewQueryBuilder[T](b.db)
}
