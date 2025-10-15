package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	accountpb "github.com/shlyyy/micro-services/api/gen"
	"github.com/shlyyy/micro-services/internal/account/model"
	"github.com/shlyyy/micro-services/pkg/db"

	"github.com/anaskhan96/go-password-encoder"

	"gorm.io/gorm"
)

type AccountService struct {
	DB *gorm.DB
	accountpb.UnimplementedAccountServiceServer
}

/*
type AccountServiceServer interface {
	GetAccountList(context.Context, *PagingRequest) (*AccountListRes, error)
	GetAccountByMobile(context.Context, *MobileRequest) (*AccountRes, error)
	GetAccountById(context.Context, *IdRequest) (*AccountRes, error)
	AddAccount(context.Context, *AddAccountRequest) (*AccountRes, error)
	UpdateAccount(context.Context, *UpdateAccountRequest) (*UpdateAccountRes, error)
	CheckPassword(context.Context, *CheckPasswordRequest) (*CheckPasswordRes, error)
	mustEmbedUnimplementedAccountServiceServer()
}
*/

func Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo == 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Model2Pb(account model.Account) *accountpb.AccountRes {
	accountRes := &accountpb.AccountRes{
		Id:       int32(account.ID),
		Mobile:   account.Mobile,
		Password: account.Password,
		Nikename: account.NikeName,
		Gender:   account.Gender,
		Role:     uint32(account.Role),
	}
	return accountRes
}

// 获取账户列表
func (s *AccountService) GetAccountList(ctx context.Context, req *accountpb.PagingRequest) (*accountpb.AccountListRes, error) {
	var accountList []model.Account
	//result := internal.DB.Find(&accountList)
	fmt.Println(req.PageNo, req.PageSize)
	result := db.GetDB().Scopes(Paginate(int(req.PageNo), int(req.PageSize))).Find(&accountList)
	if result.Error != nil {
		return nil, result.Error
	}
	accountListRes := &accountpb.AccountListRes{}

	accountListRes.Total = int32(result.RowsAffected)

	for _, account := range accountList {
		accountRes := Model2Pb(account)
		accountListRes.AccountList = append(accountListRes.AccountList, accountRes)
	}
	fmt.Println("GetAccountList Invoked...")
	return accountListRes, nil
}

// 根据 ID 获取账户信息
func (s *AccountService) GetAccountById(ctx context.Context, req *accountpb.IdRequest) (*accountpb.AccountRes, error) {
	var account model.Account
	result := db.GetDB().First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New("ACCOUNT_NOT_FOUND")
	}
	res := Model2Pb(account)
	return res, nil
}

// 根据手机号获取账户信息
func (s *AccountService) GetAccountByMobile(ctx context.Context, req *accountpb.MobileRequest) (*accountpb.AccountRes, error) {
	var account model.Account
	result := db.GetDB().Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 0 { // 没找到
		return nil, errors.New("ACCOUNT_NOT_FOUND")
	}
	res := Model2Pb(account)
	return res, nil
}

// 添加账户
func (s *AccountService) AddAccount(ctx context.Context, req *accountpb.AddAccountRequest) (*accountpb.AccountRes, error) {
	var account model.Account
	result := db.GetDB().Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 1 { // 是否存在
		return nil, errors.New("ACCOUNT_EXISTS")
	}
	// 创建账户
	account.Mobile = req.Mobile
	account.NikeName = req.NikeName
	account.Role = 1
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodePwd := password.Encode(req.Password, &options)
	account.Salt = salt
	account.Password = encodePwd
	r := db.GetDB().Create(&account)
	if r.Error != nil {
		return nil, errors.New("INTERNAL_ERROR")
	}
	accountRes := Model2Pb(account)

	return accountRes, nil
}

// 更新账户
func (s *AccountService) UpdateAccount(ctx context.Context, req *accountpb.UpdateAccountRequest) (*accountpb.UpdateAccountRes, error) {
	var account model.Account
	result := db.GetDB().First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New("ACCOUNT_NOT_FOUND")
	}

	account.Mobile = req.Mobile
	account.NikeName = req.NikeName
	account.Gender = req.Gender

	r := db.GetDB().Save(&account)
	if r.Error != nil {
		return nil, errors.New("INTERNAL_ERROR")
	}

	return &accountpb.UpdateAccountRes{Result: true}, nil
}

// 验证密码
func (s *AccountService) CheckPassword(ctx context.Context, req *accountpb.CheckPasswordRequest) (*accountpb.CheckPasswordRes, error) {
	var account model.Account
	result := db.GetDB().First(&account, req.AccountId)
	if result.Error != nil {
		return nil, errors.New("INTERNAL_ERROR")
	}
	if account.Salt == "" {
		return nil, errors.New("SALT_ERROR")
	}
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	r := password.Verify(req.Password, account.Salt, account.Password, &options)
	return &accountpb.CheckPasswordRes{Result: r}, nil
}
