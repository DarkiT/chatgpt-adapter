package cursor

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"chatgpt-adapter/core/common"
	"chatgpt-adapter/core/gin/model"
	"chatgpt-adapter/core/logger"
	"chatgpt-adapter/utils"

	"github.com/bincooo/emit.io"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/iocgo/sdk/env"
	"github.com/iocgo/sdk/stream"
)

var (
	firstDivisible bool      // 跟踪第一次能整除的情况
	lastResetDate  time.Time // 用于存储上次重置的日期
)

func fetch(ctx *gin.Context, env *env.Environment, cookie string, buffer []byte) (response *http.Response, err error) {
	count, err := checkUsage(ctx, env, 150)
	if err != nil {
		return
	}

	if count <= 0 {
		syncToken := os.Getenv("RESET_SESSION_TOKEN")
		if syncToken != "" {
			logger.Infof("系统配备了自动刷新的Token: %s ,即将自动刷新Token.", syncToken)
			err = utils.NewEmailClient().SyncSessionToken(syncToken)
			if err != nil {
				logger.Errorf("获取Token出现错误: %v", err)
				err = fmt.Errorf("当前账户Token点数不足, 自动刷新点数失败, 请添加新账号后使用！")
				return
			}
			err = fmt.Errorf("当前账户Token点数不足, 即将自动刷新, 请稍后重试当前对话！")
			return
		}
		err = fmt.Errorf("当前账户Token点数不足, 请联系管理员添加")
		return
	}

	// 检查是否是新的一天，如果是，则重置 firstDivisible
	resetFirstDivisibleIfNewDay(&firstDivisible)

	if count%50 == 0 {
		if firstDivisible {
			// 调用新封装的函数进行检查和更新
			if ok, e := checkModelAndUpdateCount(ctx, count); !ok {
				err = e
				return
			}
		} else {
			firstDivisible = true // 标记为第一次能整除
		}
	}

	logger.Infof("count: %d checksum: %s", count, ctx.GetString("checksum"))

	getApi := emit.ClientBuilder(common.HTTPClient).Context(ctx.Request.Context()).Proxies(env.GetString("server.proxied"))
	if ctx.GetBool("refresh") {
		getApi.Header("X-Amzn-Trace-Id", "Root="+uuid.New().String())
	} else {
		getApi.Header("X-Request-Id", uuid.New().String())
	}
	//response, err = emit.ClientBuilder(common.HTTPClient).
	//	Context(ctx.Request.Context()).
	//	Proxies(env.GetString("server.proxied")).
	response, err = getApi.POST("https://api2.cursor.sh/aiserver.v1.AiService/StreamChat").
		Header("Connect-Accept-Encoding", "gzip").
		Header("Authorization", "Bearer "+cookie).
		Header("Connect-Protocol-Version", "1").
		Header("Content-Type", "application/connect+proto").
		Header("User-Agent", "connect-es/1.6.1").
		Header("X-Client-Key", utils.GetClientKey()).
		Header("X-Cursor-Checksum", ctx.GetString("checksum")).
		Header("X-Cursor-Client-Version", "0.44.11").
		Header("X-Cursor-Timezone", "Asia/Macau").
		Header("X-Ghost-Mode", "true").
		Header("Connection", "keep-alive").
		Header("Transfer-Encoding", "chunked").
		Header("Host", "api2.cursor.sh").
		Bytes(buffer).
		DoC(emit.Status(http.StatusOK), emit.IsPROTO)
	return
}

func convertRequest(completion model.Completion) (buffer []byte, err error) {
	messages := stream.Map(stream.OfSlice(completion.Messages), func(message model.Keyv[interface{}]) *ChatMessage_UserMessage {
		return &ChatMessage_UserMessage{
			MessageId: uuid.NewString(),
			Role:      elseOf[int32](message.Is("role", "user"), 1, 2),
			Content:   message.GetString("content"),
		}
	}).ToSlice()
	message := &ChatMessage{
		Messages: messages,
		Instructions: &ChatMessage_Instructions{
			Instruction: "",
		},
		ProjectPath: "/path/to/project",
		Model: &ChatMessage_Model{
			Name:  completion.Model[7:],
			Empty: "",
		},
		Summary:        "",
		RequestId:      uuid.NewString(),
		ConversationId: uuid.NewString(),
	}

	protoBytes, err := proto.Marshal(message)
	if err != nil {
		return
	}

	header := int32ToBytes(0, len(protoBytes))
	buffer = append(header, protoBytes...)
	return
}

func checkUsage(ctx *gin.Context, env *env.Environment, max int) (count int, err error) {
	var (
		cookie = ctx.GetString("checktoken")
	)
	cookie, err = url.QueryUnescape(cookie)
	if err != nil {
		return
	}

	user := ""
	if strings.Contains(cookie, "::") {
		user = strings.Split(cookie, "::")[0]
	}
	response, err := emit.ClientBuilder(common.HTTPClient).
		Context(ctx.Request.Context()).
		Proxies(env.GetString("server.proxied")).
		GET("https://www.cursor.com/api/usage").
		Query("user", user).
		Header("cookie", "WorkosCursorSessionToken="+cookie).
		Header("referer", "https://www.cursor.com/settings").
		Header("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15").
		DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return
	}
	defer response.Body.Close()
	obj, err := emit.ToMap(response)
	if err != nil {
		return
	}

	for k, v := range obj {
		if !strings.Contains(k, "gpt-") {
			continue
		}
		value, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		i := value["numRequests"].(float64)
		count += int(i)
	}

	count = max - count
	return
}

/*
func genChecksum(ctx *gin.Context, env *env.Environment) string {
	token := ctx.GetString("token")
	checksum := ctx.GetHeader("x-cursor-checksum")

	if checksum == "" {
		checksum = env.GetString("cursor.checksum")
		if strings.HasPrefix(checksum, "http") {
			cacheManager := cache.CursorCacheManager()
			value, err := cacheManager.GetValue(common.CalcHex(token))
			if err != nil {
				logger.Error(err)
				return ""
			}
			if value != "" {
				return value
			}

			response, err := emit.ClientBuilder(common.HTTPClient).GET(checksum).
				DoC(emit.Status(http.StatusOK), emit.IsTEXT)
			if err != nil {
				logger.Error(err)
				return ""
			}
			checksum = emit.TextResponse(response)
			response.Body.Close()

			_ = cacheManager.SetWithExpiration(common.CalcHex(token), checksum, 30*time.Minute) // 缓存30分钟
			return checksum
		}
	}

	if checksum == "" {
		// 不采用全局设备码方式，而是用cookie产生。更换时仅需要重新抓取新的WorkosCursorSessionToken即可
		salt := strings.Split(token, ".")
		calc := func(data []byte) {
			var t byte = 165
			for i := range data {
				data[i] = (data[i] ^ t) + byte(i)
				t = data[i]
			}
		}

		// 对时间检验了
		t := time.Now()
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 30*(t.Minute()/30), 0, 0, t.Location()) // 每个半小时轮换一次
		timestamp := int64(math.Floor(float64(t.UnixMilli()) / 1e6))
		data := []byte{
			byte((timestamp >> 8) & 0xff),
			byte(timestamp & 0xff),
			byte((timestamp >> 24) & 0xff),
			byte((timestamp >> 16) & 0xff),
			byte((timestamp >> 8) & 0xff),
			byte(timestamp & 0xff),
		}
		calc(data)
		hex1 := sha256.Sum256([]byte(salt[1]))
		hex2 := sha256.Sum256([]byte(token))
		// 前面的字符生成存在问题，先硬编码
		// woc , 粗心大意呀
		checksum = fmt.Sprintf("%s%s/%s", base64.RawStdEncoding.EncodeToString(data), hex.EncodeToString(hex1[:]), hex.EncodeToString(hex2[:]))
	}
	return checksum
}
*/

func int32ToBytes(magic byte, num int) []byte {
	hex := make([]byte, 4)
	binary.BigEndian.PutUint32(hex, uint32(num))
	return append([]byte{magic}, hex...)
}

func bytesToInt32(hex []byte) int {
	return int(binary.BigEndian.Uint32(hex))
}

func elseOf[T any](condition bool, a1, a2 T) T {
	if condition {
		return a1
	}
	return a2
}

// 检查是否是新的一天，并重置 firstDivisible
func resetFirstDivisibleIfNewDay(firstDivisible *bool) {
	currentDate := time.Now().Truncate(24 * time.Hour) // 获取当前日期，不包含时间
	if lastResetDate != currentDate {
		*firstDivisible = false     // 重置为 false
		lastResetDate = currentDate // 更新最后重置日期
	}
}

// 封装函数：检查模型是否在白名单中，并更新环境变量
func checkModelAndUpdateCount(ctx *gin.Context, count int) (bool, error) {
	checkVal := fmt.Sprintf("%s_%d", time.Now().Format("2006-01-02"), count)
	envVal := os.Getenv("NOW_COUNT")

	// 直接检查白名单模型是否在列表中
	modelList := os.Getenv("WHITE_MODEL_LIST")
	if modelList != "" {
		if !inArray(ctx.GetString("modelName"), strings.Split(modelList, "|")) {
			return false, fmt.Errorf("当前账户今日高级模型点数不足, 请使用其他模型. Count: %d", count)
		}
	}

	if checkVal != envVal {
		os.Setenv("NOW_COUNT", checkVal)
	}
	return true, nil
}

func inArray(element string, array []string) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}
