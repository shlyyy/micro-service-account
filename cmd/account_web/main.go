package main

import (
	"fmt"
	"log"

	accountpb "github.com/shlyyy/micro-services/api/gen"
	"github.com/shlyyy/micro-services/pkg/config"
	"github.com/shlyyy/micro-services/pkg/logger"

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

	// 获取账户列表
	r.GET("/accounts", func(c *gin.Context) {
		pageNo := uint32(1)
		pageSize := uint32(10)
		resp, err := client.GetAccountList(c, &accountpb.PagingRequest{
			PageNo:   pageNo,
			PageSize: pageSize,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "获取账户失败"})
			return
		}
		c.JSON(200, resp)
	})

	// 登录接口（简单验证）
	r.POST("/login", func(c *gin.Context) {
		var loginReq struct {
			Mobile   string `json:"mobile"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(400, gin.H{"error": "请求格式错误"})
			return
		}

		accountResp, err := client.GetAccountByMobile(c, &accountpb.MobileRequest{
			Mobile: loginReq.Mobile,
		})
		if err != nil {
			c.JSON(404, gin.H{"error": "账户不存在"})
			return
		}

		// 验证密码
		checkResp, err := client.CheckPassword(c, &accountpb.CheckPasswordRequest{
			AccountId: (uint32)(accountResp.Id),
			Password:  loginReq.Password,
		})
		if err != nil || !checkResp.Result {
			c.JSON(401, gin.H{"error": "密码错误"})
			return
		}

		// 登录成功
		c.JSON(200, gin.H{"message": "登录成功"})
	})

	// 启动 Web 服务
	port := config.Cfg.Account.AccountWeb.Port
	r.Run(fmt.Sprintf(":%d", port))
}
