package migration

import (
	"github.com/shlyyy/micro-services/internal/account/model"
	"github.com/shlyyy/micro-services/pkg/db"
)

// 初始化 account 表
func InitAccountTable() error {
	db := db.GetDB()
	db.Migrator().DropTable(&model.Account{})
	if err := db.AutoMigrate(&model.Account{}); err != nil {
		return err
	}

	// 插入模拟数据
	accounts := []model.Account{
		{Mobile: "13800000000", Password: "password123", NikeName: "User1", Salt: "salt123", Gender: "male", Role: 1},
		{Mobile: "13800000001", Password: "password123", NikeName: "User2", Salt: "salt123", Gender: "female", Role: 2},
		{Mobile: "13800000002", Password: "password123", NikeName: "User3", Salt: "salt123", Gender: "male", Role: 1},
		{Mobile: "13800000003", Password: "password123", NikeName: "User4", Salt: "salt123", Gender: "female", Role: 2},
	}
	if err := db.Create(&accounts).Error; err != nil {
		return err
	}
	return nil
}
