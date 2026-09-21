package query

import (
	"gorm.io/gorm"
)

// ---------------- 泛型查询构造器 ----------------
// QueryBuilder[T] 泛型构造器，T 为模型结构体（非指针）
// QueryBuilder 临时查询构造器，只在每次查询时创建，存放链式会话，带Eq等方法
type QueryBuilder[T any] struct {
	db *gorm.DB
}

// NewQueryBuilder 创建泛型构造器，绑定模型 T
func NewQueryBuilder[T any](db *gorm.DB) *QueryBuilder[T] {
	var t T
	return &QueryBuilder[T]{
		db: db.Model(&t),
	}
}

// Where 条件查询
func (qb *QueryBuilder[T]) Where(query interface{}, args ...interface{}) *QueryBuilder[T] {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// OrderAsc 升序排序
func (qb *QueryBuilder[T]) OrderAsc(field string) *QueryBuilder[T] {
	qb.db = qb.db.Order(field + " ASC")
	return qb
}

// OrderDesc 降序排序
func (qb *QueryBuilder[T]) OrderDesc(field string) *QueryBuilder[T] {
	qb.db = qb.db.Order(field + " DESC")
	return qb
}

// Eq 相等条件
func (qb *QueryBuilder[T]) Eq(field string, val any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" = ?", val)
	return qb
}

// EqIf 当 ok == true 时才追加等值条件
func (qb *QueryBuilder[T]) EqIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Eq(field, val)
	}
	return qb
}

// Neq !=
func (qb *QueryBuilder[T]) Neq(field string, val any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" != ?", val)
	return qb
}

// NeqIf 不等，条件成立才加
func (qb *QueryBuilder[T]) NeqIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Neq(field, val)
	}
	return qb
}

// Gt >
func (qb *QueryBuilder[T]) Gt(field string, val any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" > ?", val)
	return qb
}

// GtIf 大于
func (qb *QueryBuilder[T]) GtIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Gt(field, val)
	}
	return qb
}

// Lt <
func (qb *QueryBuilder[T]) Lt(field string, val any) *QueryBuilder[T] {
	qb = qb.Lt(field, val)
	return qb
}

// LtIf 小于
func (qb *QueryBuilder[T]) LtIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Lt(field, val)
	}
	return qb
}

// Gte >=
func (qb *QueryBuilder[T]) Gte(field string, val any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" >= ?", val)
	return qb
}

// GteIf 大于等于
func (qb *QueryBuilder[T]) GteIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Gte(field, val)
	}
	return qb
}

// Lte <=
func (qb *QueryBuilder[T]) Lte(field string, val any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" <= ?", val)
	return qb
}

// LteIf 小于等于
func (qb *QueryBuilder[T]) LteIf(ok bool, field string, val any) *QueryBuilder[T] {
	if ok {
		qb = qb.Lte(field, val)
	}
	return qb
}

// In IN 查询
func (qb *QueryBuilder[T]) In(field string, vals []any) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" IN ?", vals)
	return qb
}

// InIf in 查询，vals 非nil/非空才生效
func (qb *QueryBuilder[T]) InIf(ok bool, field string, vals []any) *QueryBuilder[T] {
	if ok {
		qb = qb.In(field, vals)
	}
	return qb
}

// Like 模糊查询
func (qb *QueryBuilder[T]) Like(field string, val string) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" LIKE ?", "%"+val+"%")
	return qb
}

// LikeIf 模糊查询
func (qb *QueryBuilder[T]) LikeIf(ok bool, field string, val string) *QueryBuilder[T] {
	if ok {
		qb = qb.Like(field, val)
	}
	return qb
}

// Order 排序
func (qb *QueryBuilder[T]) Order(expr string) *QueryBuilder[T] {
	qb.db = qb.db.Order(expr)
	return qb
}

// Limit 限制行数
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 偏移
func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.db = qb.db.Offset(offset)
	return qb
}

// GetDB 返回原生 *gorm.DB，兜底扩展
func (qb *QueryBuilder[T]) GetDB() *gorm.DB {
	return qb.db
}

// // First 查询单条
// func (qb *QueryBuilder[T]) First(dest interface{}, conds ...interface{}) *gorm.DB {
// 	return qb.db.First(dest, conds...)
// }

// // Find 查询列表
// func (qb *QueryBuilder[T]) Find(dest interface{}, conds ...interface{}) *gorm.DB {
// 	return qb.db.Find(dest, conds...)
// }

// First 查询单条
func (qb *QueryBuilder[T]) First() (*T, error) {
	var entity T
	err := qb.db.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Find 查询列表
func (qb *QueryBuilder[T]) Find() ([]T, error) {
	var list []T
	err := qb.db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Count 统计总数
func (qb *QueryBuilder[T]) Count() (int64, error) {
	var cnt int64
	err := qb.db.Count(&cnt).Error
	return cnt, err
}
