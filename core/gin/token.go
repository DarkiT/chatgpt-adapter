package gin

import (
	"net/http"

	"chatgpt-adapter/core/logger"
	"chatgpt-adapter/utils"

	"github.com/gin-gonic/gin"
)

// @POST(path = "/api/token")
func (h *Handler) token(c *gin.Context) {
	var req *utils.ConfigRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := gin.H{
		"status":  "success",
		"message": "配置保存失败",
	}
	if ok, req := utils.AuthToken(req, true); ok {
		logger.Info("配置保存成功")
		message = gin.H{
			"status":  "success",
			"message": "配置保存成功",
			"data":    req,
		}
	}
	c.JSON(http.StatusOK, message)
}

// @POST(path = "/api/config")
func (h *Handler)  config(c *gin.Context) {
	c.JSON(200, gin.H{
		"needCode":         false,
		"hideUserApiKey":   false,
		"disableGPT4":      false,
		"hideBalanceQuery": true,
		"disableFastLink":  false,
		"customModels":     "",
		"defaultModel":     "deepseek-r1",
	})
}
