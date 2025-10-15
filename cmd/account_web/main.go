package main

import (
	"fmt"
	"log"

	"github.com/shlyyy/micro-service-account/api/accountpb"
	"github.com/shlyyy/micro-service-account/internal/account/handler"
	"github.com/shlyyy/micro-service-account/pkg/config"
	"github.com/shlyyy/micro-service-account/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		panic(fmt.Errorf("加载配置失败: %w", err))
	}

	// 初始化日志
	logger.InitLogger(config.Cfg.Logger)

	// 连接 gRPC 服务
	serverPort := config.Cfg.Account.AccountServer.GrpcServerPort
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到 gRPC 服务: %v", err)
	}
	defer conn.Close()

	client := accountpb.NewAccountServiceClient(conn)

	// 创建 Gin 路由
	r := gin.Default()

	// 初始化 Web 路由处理
	handler.NewAccountHandler(r, client)

	// 启动 Web 服务
	port := config.Cfg.Account.AccountWeb.Port
	r.Run(fmt.Sprintf(":%d", port))
}
