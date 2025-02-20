package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darkit/slog"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	hostname string
	workPath = filepath.Dir(os.Args[0])
	configDb = filepath.Join(workPath, fmt.Sprintf("%s.db", strings.ReplaceAll(filepath.Base(os.Args[0]), ".exe", "")))
)

func init() {
	hostname, _ = os.Hostname()
	hostname = hex.EncodeToString(hmac.New(sha1.New, []byte(hostname)).Sum([]byte(hostname)))
}

type ConfigRequest struct {
	ID       string `json:"id,omitempty" form:"id"`
	Cookie   string `json:"cookie" form:"cookie" binding:"required"`
	Checksum string `json:"checksum" form:"checksum"`
}

type jwtInfo struct {
	jwt.RegisteredClaims
}

func GetAuthToken() *ConfigRequest {
	_, tk := AuthToken(nil, false)
	// client := NewEmailClient()
	// if ok {
	// 	syncToken := os.Getenv("RESET_SESSION_TOKEN")
	// 	if syncToken != "" {
	// 		userID, cookie, ok := ValidateToken(tk.Cookie)
	// 		if !ok {
	// 			return nil
	// 		}
	// 		usage, err := client.GetUsage(userID, cookie)
	// 		if err == nil {
	// 			if usage.GPT4.NumRequests >= 150 {
	// 				_ = client.SyncSessionToken(syncToken)
	// 			}
	// 		} else {
	// 			slog.Error("获取usage失败", "error", err.Error())
	// 		}
	// 	}
	// }

	return tk
}

func GetUsage() *GPTUsageData {
	ok, tk := AuthToken(nil, false)
	client := NewEmailClient()
	if ok {

		userID, cookie, ok := ValidateToken(tk.Cookie)
		if !ok {
			return nil
		}
		usage, err := client.GetUsage(userID, cookie)
		if err == nil {
			return usage
		} else {
			slog.Error("获取usage失败", "error", err.Error())
		}
	}

	return nil
}

func AuthToken(token *ConfigRequest, isEncode bool) (bool, *ConfigRequest) {
	if err := os.Chdir(workPath); err != nil {
		slog.Errorf("切换WorkPath失败: %s", workPath)
		return false, nil
	}
	if isEncode {
		if token.Checksum == "" {
			token.Checksum, _ = GenChecksum(token)
		}
		var buf strings.Builder
		if err := gob.NewEncoder(&buf).Encode(token); err != nil {
			slog.Error("编码Token文件失败!", "error", err)
			return false, nil
		}
		if err := os.WriteFile(configDb, []byte(buf.String()), 0o600); err != nil {
			slog.Error("设置Token文件失败!", "error", err)
			return false, nil
		}
	} else {
		db, err := os.ReadFile(configDb)
		if err != nil {
			return false, nil
		}
		if err = gob.NewDecoder(bytes.NewReader(db)).Decode(&token); err != nil {
			slog.Error("解码Token文件失败", "error", err)
			return false, nil
		}
	}

	return true, token
}

func GenChecksum(req *ConfigRequest) (string, bool) {
	var checksum string

	if req.Cookie == "" {
		return "", false
	}
	salt := strings.Split(req.Cookie, ".")

	// 时间处理，固定为当天的零点
	t := time.Now()
	//t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // 每天的零点时间
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 17*(t.Minute()/17), 0, 0, t.Location()) // 每17分钟轮换一次
	timestamp := int64(math.Floor(float64(t.UnixMilli()) / 1e6))                                  // 毫秒级时间戳，单位为秒

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
	hex1 := sha256.Sum256([]byte(salt[1])) // 基于salt的sha256
	checksum = fmt.Sprintf("%s%s/%s", base64.RawStdEncoding.EncodeToString(data), hex.EncodeToString(hex1[:]), uuid.NewSHA1(uuid.NameSpaceDNS, hmac.New(sha1.New, []byte(req.Cookie)).Sum([]byte(req.Cookie))).String())

	if checksum != req.Checksum {
		req.Checksum = checksum
		return checksum, true
	}

	return checksum, false
}

func ValidateToken(tokenString string) (string, string, bool) {
	if strings.HasPrefix(tokenString, "user_") {
		decodedURL, err := url.QueryUnescape(tokenString)
		if err != nil {
			return "", "", false
		}
		tks := strings.Split(decodedURL, "::")
		if len(tks) != 2 {
			return "", "", false
		}
		tokenString = tks[1]
	}
	token, _ := jwt.ParseWithClaims(tokenString, &jwtInfo{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if token == nil {
		return "", "", false
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", "", false
	}
	userID := strings.Split(subject, "|")
	if len(userID) != 2 {
		return "", "", false
	}
	return userID[1], tokenString, true
}

func GetClientKey() string {
	return hostname
}
