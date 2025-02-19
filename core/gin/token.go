package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		"code":    500,
		"status":  "error",
		"message": "配置保存失败",
	}

	if req.Cookie != "" {
		_, cookie, ok := utils.ValidateToken(req.Cookie)
		if ok {
			req.Cookie = cookie
			if ok, req = utils.AuthToken(req, true); ok {
				message = gin.H{
					"code":    200,
					"status":  "success",
					"message": "配置保存成功",
					"data":    req,
				}
			}
		} else {
			client := utils.NewEmailClient()
			err := client.LoginDeepControl(req.Cookie)
			if err != nil {
				message = gin.H{
					"code":    500,
					"status":  "error",
					"message": err.Error(),
				}
			} else {
				message = gin.H{
					"code":    200,
					"status":  "success",
					"message": "登录成功",
				}
			}
		}

	}

	c.JSON(http.StatusOK, message)
}

// @POST(path = "/api/config")
func (h *Handler) config(c *gin.Context) {
	c.JSON(200, gin.H{
		"needCode":         false,
		"hideUserApiKey":   false,
		"disableGPT4":      false,
		"hideBalanceQuery": true,
		"disableFastLink":  false,
		"customModels":     "",
		"defaultModel":     "zishuo/deepseek-r1",
	})
}

// Notification 结构体用于表示通知
type Notification struct {
	Type  string `json:"type"`            // 通知类型
	Email string `json:"email,omitempty"` // 邮箱地址
	Code  string `json:"code,omitempty"`  // 验证码内容
	Usage any    `json:"usage,omitempty"` // 套餐余量信息
}

// @GET(path = "/api/notifications")
func (h *Handler) notifications(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	client := utils.NewEmailClient()
	emailAddr, err := client.GetEmailAddress()
	if err != nil {
		sendSSEError(c, "无法获取邮箱地址")
		return
	}

	// 先发用户信息套餐信息（仅一次）
	sendNotification(c, Notification{
		Type:  "usage",
		Usage: utils.GetUsage(),
	})

	// 先发送邮箱地址（仅一次）
	sendNotification(c, Notification{
		Type:  "email",
		Email: emailAddr.Address,
	})

	// 轮询验证码
	for {
		code, err := client.GetVerificationCode()
		if err != nil {
			//sendSSEError(c, "验证码获取失败,请稍后重试...")
			sendSSEError(c, err.Error())
			return
		}
		if code != "" {
			sendNotification(c, Notification{
				Type: "code",
				Code: code,
			})
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// @Any(path = "/api/cache/")
func (h *Handler) cache(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func sendNotification(c *gin.Context, n Notification) {
	data, _ := json.Marshal(n)
	fmt.Fprintf(c.Writer, "event: message\ndata: %s\n\n", data)
	c.Writer.Flush()
}

func sendSSEError(c *gin.Context, msg string) {
	fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", msg)
	c.Writer.Flush()
}
