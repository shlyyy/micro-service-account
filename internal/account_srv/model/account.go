package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Mobile   string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(64);not null"`
	NikeName string `gorm:"type:varchar(32)"`
	Salt     string `gorm:"type:varchar(16)"`
	Gender   string `gorm:"type:varchar(6);default:male"`
	Role     int    `gorm:"type:int;default:1;comment:'1-普通用户,2-管理员'"`
}
