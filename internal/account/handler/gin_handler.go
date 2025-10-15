package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shlyyy/micro-services/api/accountpb"
)

func NewAccountHandler(r *gin.Engine, client accountpb.AccountServiceClient) {
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
			AccountId: uint32(accountResp.Id),
			Password:  loginReq.Password,
		})
		if err != nil || !checkResp.Result {
			c.JSON(401, gin.H{"error": "密码错误"})
			return
		}

		// 登录成功
		c.JSON(200, gin.H{"message": "登录成功"})
	})
}
