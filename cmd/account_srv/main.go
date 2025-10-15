package main

import (
	"fmt"
	"net"

	"github.com/shlyyy/micro-service-account/api/accountpb"
	"github.com/shlyyy/micro-service-account/internal/account/migration"
	"github.com/shlyyy/micro-service-account/internal/account/service"
	"github.com/shlyyy/micro-service-account/pkg/config"
	"github.com/shlyyy/micro-service-account/pkg/db"
	"github.com/shlyyy/micro-service-account/pkg/logger"
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

	if err := migration.InitAccountTable(); err != nil {
		logger.Error("初始化 account 表失败:", err)
		return
	}
	logger.Info("account 表初始化成功")

	// 初始化服务
	accountService := &service.AccountService{DB: db.GetDB()}
	grpcServer := grpc.NewServer()

	accountpb.RegisterAccountServiceServer(grpcServer, accountService)

	// 启动 gRPC 服务器
	port := config.Cfg.Account.AccountServer.GrpcServerPort
	logger.Infof("启动 gRPC 服务器，监听端口: %d", port)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Errorf("无法监听端口: %v", err)
		return
	}

	logger.Info("account_srv 服务已启动")

	if err := grpcServer.Serve(listen); err != nil {
		logger.Errorf("启动 gRPC 服务失败: %v", err)
	}
}
