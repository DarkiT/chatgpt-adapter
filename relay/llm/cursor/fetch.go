package cursor

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"chatgpt-adapter/core/cache"
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
	g_checksum = ""
)

func fetch(ctx *gin.Context, env *env.Environment, cookie string, buffer []byte) (response *http.Response, err error) {
	count, err := checkUsage(ctx, env, 150)
	if err != nil {
		return
	}
	if count <= 0 {
		err = fmt.Errorf("invalid usage")
		return
	}

	//accept-encoding: gzip
	//authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHx1c2VyXzAxSk0xNDIwRDVBR0UxMEU1UjJEMVE0S0UzIiwidGltZSI6IjE3Mzk0OTk2NDEiLCJyYW5kb21uZXNzIjoiOWZjYTUzMDgtZGU3My00YzliIiwiZXhwIjo0MzMxNDk5NjQxLCJpc3MiOiJodHRwczovL2F1dGhlbnRpY2F0aW9uLmN1cnNvci5zaCIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJhdWQiOiJodHRwczovL2N1cnNvci5jb20ifQ.v9LhtzlZgm7uMt9XYFJ8JCz8ajjhKmI10w8xW1I3mxc
	//connect-protocol-version: 1
	//content-type: application/proto
	//cookie:
	//traceparent: 00-ef9953bc56d444608e24aa54d31e9a1c-8c5cfa2d832f5d0b-00
	//user-agent: connect-es/1.6.1
	//x-client-key: 4fc521e53b3204feafb9d639c3581b4ad6d01f64459a5200a3ce572404891e56
	//x-cursor-checksum: LpyehxCq565f478973be4748a6523d4b3090ef2250cff6c8bfda465aaa4ef0aafb74e52d/0bcdac88-d1e1-41c0-aa7b-73105ca1628f
	//x-cursor-client-version: 0.44.11
	//x-cursor-timezone: Asia/Macau
	//x-ghost-mode: true
	//host: api2.cursor.sh
	//connection: keep-alive
	//transfer-encoding: chunked

	var authToken *utils.ConfigRequest

	if cookie == "" {
		authToken = utils.GetAuthToken(false)
	} else {
		authToken.Cookie = cookie
	}

	if authToken.Checksum == "" {
		authToken.Checksum = utils.GenChecksum(authToken)
	}

	if authToken.Checksum == "" {
		authToken.Checksum = genChecksum(ctx, env)
	}

	logger.Debugf("checksum: %s", authToken.Checksum)

	response, err = emit.ClientBuilder(common.HTTPClient).
		Context(ctx.Request.Context()).
		Proxies(env.GetString("server.proxied")).
		POST("https://api2.cursor.sh/aiserver.v1.AiService/StreamChat").
		//Header("authorization", "Bearer "+cookie).
		//Header("content-type", "application/connect+proto").
		//Header("connect-accept-encoding", "gzip,br").
		//Header("connect-protocol-version", "1").
		//Header("user-agent", "connect-es/1.4.0").
		//Header("x-cursor-checksum", genChecksum(ctx, env)).
		//Header("x-cursor-client-version", "0.42.3").
		//Header("x-cursor-timezone", "Asia/Shanghai").
		//Header("host", "api2.cursor.sh").
		//Header("X-Amzn-Trace-Id", "Root="+uuid.New().String())
		Header("Connect-Accept-Encoding", "gzip").
		Header("Authorization", "Bearer "+authToken.Cookie).
		Header("Connect-Protocol-Version", "1").
		Header("Content-Type", "application/connect+proto").
		Header("User-Agent", "connect-es/1.6.1").
		Header("X-Client-Key", utils.GetClientKey()).
		Header("X-Cursor-Checksum", authToken.Checksum).
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
		cookie = ctx.GetString("token")
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
		hex1 := sha256.Sum256([]byte(salt[1])) // 基于salt的sha256
		hex2 := sha256.Sum256([]byte(token))   // 基于cookie的sha256

		// 拼接最终的checksum，包含时间戳、hex1和hex2  //hex.EncodeToString(hex2[:])
		checksum = fmt.Sprintf("%s%s/%s", base64.RawStdEncoding.EncodeToString(data), hex.EncodeToString(hex1[:]), uuid.NewSHA1(uuid.NameSpaceDNS, hmac.New(sha1.New, hex2[:]).Sum(hex2[:])).String())
	}
	return checksum
}

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
