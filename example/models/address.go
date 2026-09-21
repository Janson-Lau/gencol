package models

import (
	"github.com/Janson-Lau/gencol/example/meta"
	"github.com/Janson-Lau/gencol/query"
	"gorm.io/gorm"
)

type Address struct {
	Id        int64  `gorm:"column:id;type:bigint;primary_key;comment:id"`
	Address   string `gorm:"column:address;type:varchar(100);not null;comment:地址"`
	IsActive  bool   `gorm:"column:is_active;type:tinyint;not null;default:1;comment:是否启用"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;autoCreateTime:milli;comment:创建时间"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;autoUpdateTime:milli;comment:更新时间"`
}

type AddressModel struct {
	db *gorm.DB
}

func NewAddressModel(db *gorm.DB) *AddressModel {
	return &AddressModel{db: db}
}

// NewQuery 创建泛型查询器，绑定 Address
func (p *AddressModel) NewQuery() *query.QueryBuilder[Address] {
	return query.NewQueryBuilder[Address](p.db)
}

// GetByAddress 你的目标方法，链式 Eq
func (p *AddressModel) GetByAddress(address string) (*Address, error) {
	return p.NewQuery().
		Eq(meta.Address.Address, address).
		Eq(meta.Address.IsActive, true).
		First()
}
