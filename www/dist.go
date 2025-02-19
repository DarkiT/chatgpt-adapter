package www

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"chatgpt-adapter/core/logger"
	"github.com/darkit/machineid"
	"github.com/darkit/machineid/cert"
)

//go:embed all:public
var Dist embed.FS

var (
	auth, _  = cert.New()
	authFile = fmt.Sprintf("%s.crt", filepath.Base(os.Args[0]))
)

func init() {
	protectedID, err := machineid.ProtectedID("CursorTools")
	if err != nil {
		logger.Fatal("获取设备机器码出错")
	}

	// 首次检查授权
	if err = validateAuth(protectedID); err != nil {
		handleAuthError(protectedID, err)
		return
	}

}

// validateAuth 验证授权
func validateAuth(protectedID string) error {
	// 读取授权文件
	certData, err := os.ReadFile(authFile)
	if err != nil {
		return fmt.Errorf("读取授权文件失败, 请检查工作目录是否有 %s 文件", authFile)
	}

	// 验证授权
	if err = auth.ValidateCert(certData, protectedID); err != nil {
		return fmt.Errorf("请使用合法授权文件")
	}

	return nil
}

// handleAuthError 处理授权错误
func handleAuthError(protectedID string, err error) {
	logger.Warnf("您的设备机器码为: %s", protectedID)
	logger.Errorf("授权校验不通过: %v", err)

	// 进度条倒计时
	for i := 0; i <= 50; i++ {
		// 构建进度条：左边文本 + 进度条 + 右边百分比
		progressBar := fmt.Sprintf("%s [W] 程序即将退出: [%s>%s] %3d%%",
			time.Now().Format("2006-01-02 15:04:05.000"),
			strings.Repeat("=", i),    // 已完成部分
			strings.Repeat(" ", 50-i), // 未完成部分
			i*2)
		// 使用\r回到行首，刷新进度条
		fmt.Printf("\r%s", progressBar)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println() // 最后换行

	os.Exit(1)
}
