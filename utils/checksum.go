package utils

import (
	"bufio"
	"bytes"
	"chatgpt-adapter/core/logger"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	hostname string
	workPath = filepath.Dir(os.Args[0])
	//cache    = freecache.NewCache(512 * 1024) // 0.5 MB的字节数
	configDb = filepath.Join(workPath, fmt.Sprintf("%s.db", strings.ReplaceAll(filepath.Base(os.Args[0]), ".exe", "")))
)

type ConfigRequest struct {
	Cookie   string `json:"cookie" form:"cookie" binding:"required"`
	Checksum string `json:"checksum" form:"checksum"`
}

func init() {
	hostname, _ = os.Hostname()
	hostname = hex.EncodeToString(hmac.New(sha1.New, []byte(hostname)).Sum([]byte(hostname)))
}

func GetAuthToken(wait bool) *ConfigRequest {
	ok, tk := AuthToken(nil, false)
	if !ok && wait {
		client := NewEmailClient()
		logger.Warn("读取Token文件失败, 即将启动快速注册流程.")
		emailAddr, err := client.GetEmailAddress()
		if err != nil {
			logger.Errorf("邮件客户端初始化失败: %v", err.Error())
			logger.Error("请访问 www.cursor.com 手动完成账号注册.")
			logger.Error("然后按F12打开开发者工具, 找到“应用-Cookies”中名为 WorkosCursorSessionToken 的值并保存至当前目录下的cookie.txt中.")
			return nil
		}
		logger.Infof("成功获取邮箱地址: %s 请访问 www.cursor.com 并使用该邮箱完成账号注册", emailAddr.Address)
		logger.Info("开始等待注册验证码.")
		// 轮询收件箱
		var emails []Email
		for attempt := 1; attempt <= maxPollAttempts; attempt++ {
			emails, err = client.PollInbox()
			if err != nil {
				logger.Errorf("轮询失败（尝试次数 %d）: %v", attempt, err)
				continue
			}

			if len(emails) > 0 {
				logger.Infof("第 %d 次尝试成功收到邮件", attempt)
				break
			}

			time.Sleep(pollInterval)
		}

		if len(emails) == 0 {
			logger.Error("超过最大轮询次数仍未收到邮件, 请访问 www.cursor.com 手动完成账号注册.")
			logger.Error("然后按F12打开开发者工具, 找到“应用-Cookies”中名为 WorkosCursorSessionToken 的值并保存至当前目录下的cookie.txt中.")
			return nil
		}

		// 输出邮件内容
		var code string
		for _, email := range emails {
			code, err = extractVerificationCode(email.Body)
			if err == nil {
				logger.Infof("提取到验证码 %s 请录入验证码完成注册验证. ", code)
				logger.Info("然后进入管理页面后打开开发者工具, 找到“应用-Cookies”中名为 WorkosCursorSessionToken 的值复制并发送给我.")
				break
			}
		}
		if code != "" {
			time.Sleep(3 * time.Second)
			_ = OpenUrl("https://authenticator.cursor.sh/?redirect_uri=https%3A%2F%2Fcursor.com%2Fapi%2Fauth%2Fcallback")

			reader := bufio.NewReader(os.Stdin)
			session := promptInput(reader, "请输入获取到的WorkosCursorSessionToken内容: ")
			if session != "" && validateString(session) {
				ok, tk = AuthToken(&ConfigRequest{Cookie: session}, true)
				if !ok {
					logger.Error("保存认证文件失败!")
					logger.Warn("请将获取到的WorkosCursorSessionToken内容手动保存至当前目录下的cookie.txt中.")
				}
				return tk
			}
		}
		logger.Error("读取Token文件失败,请访问 www.cursor.com 完成账号注册.")
		logger.Error("然后按F12打开开发者工具, 找到“应用-Cookies”中名为 WorkosCursorSessionToken 的值并保存至当前目录下的cookie.txt中.")
		return nil
	}
	return tk
}

func AuthToken(token *ConfigRequest, isEncode bool) (bool, *ConfigRequest) {
	if err := os.Chdir(workPath); err != nil {
		logger.Errorf("切换WorkPath失败: %s", workPath)
		return false, nil
	}
	if isEncode {
		if token.Checksum == "" {
			token.Checksum = GenChecksum(token)
		}
		var buf strings.Builder
		if err := gob.NewEncoder(&buf).Encode(token); err != nil {
			logger.Error("编码Token文件失败!", "error", err)
			return false, nil
		}
		if err := os.WriteFile(configDb, []byte(buf.String()), 0o600); err != nil {
			logger.Error("设置Token文件失败!", "error", err)
			return false, nil
		}
	} else {
		db, err := os.ReadFile(configDb)
		if err != nil {
			logger.Error("读取Token文件失败!", "error", err)
			return false, nil
		}
		if err = gob.NewDecoder(bytes.NewReader(db)).Decode(&token); err != nil {
			logger.Error("解码Token文件失败", "error", err)
			return false, nil
		}
	}

	return true, token
}

func GenChecksum(req *ConfigRequest) string {
	//const (
	//	checkSumKey   = "checksum"   // Checksum缓存使用的key
	//	expireSeconds = 24 * 60 * 60 // 缓存有效期为1天（86400秒）
	//)
	var checksum string
	//if cachedValue, err := cache.Get([]byte(checkSumKey)); err == nil {
	//	return string(cachedValue)
	//}

	if req.Checksum == "" {
		salt := strings.Split(req.Cookie, ".")

		// 时间处理，固定为当天的零点
		t := time.Now()
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // 每天的零点时间
		timestamp := int64(math.Floor(float64(t.UnixMilli()) / 1e6))          // 毫秒级时间戳，单位为秒

		// 用timestamp构造一个固定长度的data数据，保证每天相同
		data := []byte{
			byte((timestamp >> 8) & 0xff),
			byte(timestamp & 0xff),
			byte((timestamp >> 24) & 0xff),
			byte((timestamp >> 16) & 0xff),
			byte((timestamp >> 8) & 0xff),
			byte(timestamp & 0xff),
		}

		// 使用calc函数处理data内容
		calc := func(data []byte) {
			var t byte = 165
			for i := range data {
				data[i] = (data[i] ^ t) + byte(i)
				t = data[i]
			}
		}
		calc(data)

		// 计算盐值和cookie的sha256哈希
		hex1 := sha256.Sum256([]byte(salt[1]))    // 基于salt的sha256
		hex2 := sha256.Sum256([]byte(req.Cookie)) // 基于cookie的sha256

		// 拼接最终的checksum，包含时间戳、hex1和hex2  //hex.EncodeToString(hex2[:])
		checksum = fmt.Sprintf("%s%s/%s", base64.RawStdEncoding.EncodeToString(data), hex.EncodeToString(hex1[:]), uuid.NewSHA1(uuid.NameSpaceDNS, hmac.New(sha1.New, hex2[:]).Sum(hex2[:])).String())
	}

	// 将获取到的checksum存入缓存，并设置有效期为1天
	//_ = cache.Set([]byte(checkSumKey), []byte(checksum), expireSeconds)

	return checksum
}

func GetClientKey() string {
	return hostname
}

// promptInput 获取用户输入
func promptInput(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("输入错误，请重试")
			continue
		}

		// 清理输入
		input = strings.TrimSpace(input)
		if input != "" {
			return strings.Trim(input, " ")
		}
		fmt.Println("输入不能为空，请重试")
	}
}

// 验证字符串是否符合预期格式
func validateString(str string) bool {
	// 检查字符串是否以 "user_" 开头
	if !strings.HasPrefix(str, "user_") {
		return false
	}

	// 检查字符串长度，假设字符串长度需要在一定范围内
	if len(str) < 250 || len(str) > 600 {
		return false
	}

	// 检查字符串是否包含点，作为JWT的标志
	if !strings.Contains(str, ".") {
		return false
	}

	// 使用正则表达式检查字符串是否符合预期格式
	// 假设JWT包含三个部分：header、payload和signature，每部分用点分隔
	re := regexp.MustCompile(`^[a-zA-Z0-9-_]+(?:\.[a-zA-Z0-9-_]+){2}\.[a-zA-Z0-9-_]+$`)
	if !re.MatchString(str) {
		return false
	}

	// 可选：检查是否包含URL编码的字符
	if strings.Contains(str, "%") {
		// 检查是否包含URL编码的字符，例如%3A
		re = regexp.MustCompile(`%[0-9A-Fa-f]{2}`)
		if !re.MatchString(str) {
			return false
		}
	}

	return true
}
