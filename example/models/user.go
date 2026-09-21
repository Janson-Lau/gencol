package models

import (
	"github.com/Janson-Lau/gencol/example/meta"
	"github.com/Janson-Lau/gencol/query"
	"gorm.io/gorm"
)

type User struct {
	Id        int64  `gorm:"column:id;type:bigint;primary_key;comment:id"`
	Username  string `gorm:"column:username;type:varchar(100);not null;comment:名称"`
	Password  string `gorm:"column:password;type:varchar(100);not null;comment:密码"`
	Phone     string `gorm:"column:phone;type:varchar(100);not null;comment:手机"`
	Email     string `gorm:"column:email;type:varchar(100);not null;comment:邮箱"`
	IsActive  bool   `gorm:"column:is_active;type:tinyint;not null;default:1;comment:是否启用"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;autoCreateTime:milli;comment:创建时间"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;autoUpdateTime:milli;comment:更新时间"`
}

type UserModel struct {
	*query.BaseModel[User]
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{BaseModel: query.NewBaseModel[User](db)}
}

// GetByUsername 你的目标方法，链式 Eq
func (p *UserModel) GetByUsername(username string) (*User, error) {
	return p.NewQuery().
		Eq(meta.User.Username, username).
		Eq(meta.User.IsActive, true).
		First()
}

// ListUser 示例：列表+分页
func (p *UserModel) ListUser(keyword string, page, size int) ([]User, int64, error) {
	q := p.NewQuery().
		Eq(meta.User.IsActive, true).
		Like(meta.User.Username, keyword)

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	list, err := q.Offset((page - 1) * size).Limit(size).Find()
	return list, total, err
}
