package main

import (
	"fmt"
	"net"

	accountpb "github.com/shlyyy/micro-services/api/gen"
	"github.com/shlyyy/micro-services/internal/account/model"
	"github.com/shlyyy/micro-services/internal/account/service"
	"github.com/shlyyy/micro-services/pkg/config"
	"github.com/shlyyy/micro-services/pkg/db"
	"github.com/shlyyy/micro-services/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		panic(fmt.Errorf("加载配置失败: %w", err))
	}

	// 初始化日志
	logger.InitLogger(config.Cfg.Logger)

	// 数据库连接
	if err := db.InitDB(&config.Cfg.Database); err != nil {
		logger.Error("数据库初始化失败:", err)
		return
	}
	logger.Info("数据库初始化成功")

	if err := initAccountTable(); err != nil {
		logger.Error("初始化 account 表失败:", err)
		return
	}
	logger.Info("account 表初始化成功")

	// 初始化服务
	accountService := &service.AccountService{DB: db.GetDB()}

	// 启动 gRPC 服务器
	port := config.Cfg.Account.AccountServer.GrpcServerPort
	logger.Infof("启动 gRPC 服务器，监听端口: %d", port)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Errorf("无法监听端口: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	accountpb.RegisterAccountServiceServer(grpcServer, accountService)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Errorf("启动 gRPC 服务失败: %v", err)
	}

	logger.Info("account_srv 服务已启动")

	select {}
}

// 初始化 account 表
func initAccountTable() error {
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
