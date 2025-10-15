package migration

import (
	"crypto/md5"
	"fmt"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/shlyyy/micro-service-account/internal/account/model"
	"github.com/shlyyy/micro-service-account/pkg/db"
)

// 初始化 account 表
func InitAccountTable() error {
	sqldb := db.GetDB()
	sqldb.Migrator().DropTable(&model.Account{})
	if err := sqldb.AutoMigrate(&model.Account{}); err != nil {
		return err
	}

	// 预先添加几个用户
	pwd := "123456"
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodePwd := password.Encode(pwd, &options)

	for i := 1; i <= 5; i++ {
		account := model.Account{
			Mobile:   "1380000000" + fmt.Sprintf("%d", i),
			Nickname: "user" + fmt.Sprintf("%d", i),
			Password: encodePwd,
			Salt:     salt,
			Role:     i%2 + 1,
		}
		sqldb.Create(&account)
	}
	return nil
}
