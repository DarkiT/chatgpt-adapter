package pg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bincooo/chatgpt-adapter/v2/internal/common"
	"github.com/bincooo/chatgpt-adapter/v2/internal/middle"
	"github.com/bincooo/chatgpt-adapter/v2/pkg"
	"github.com/bincooo/chatgpt-adapter/v2/pkg/gpt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bincooo/sdio"
)

type modelPayload struct {
	BatchId                 string  `json:"batchId"`
	CfgScale                int32   `json:"cfg_scale"`
	DreamBoothModel         string  `json:"dream_booth_model"`
	Filter                  string  `json:"filter"`
	GenerateVariants        bool    `json:"generateVariants"`
	GuidanceScale           int32   `json:"guidance_scale"`
	Width                   int32   `json:"width"`
	Height                  int32   `json:"height"`
	HighNoiseFrac           float32 `json:"high_noise_frac"`
	InitImageFromPlayground bool    `json:"initImageFromPlayground"`
	IsPrivate               bool    `json:"isPrivate"`
	ModelType               string  `json:"modelType"`
	NegativePrompt          string  `json:"negativePrompt"`
	NumImages               int32   `json:"num_images"`
	Prompt                  string  `json:"prompt"`
	Sampler                 int32   `json:"sampler"`
	Seed                    int32   `json:"seed"`
	StatusUUID              string  `json:"statusUUID"`
	Steps                   int32   `json:"steps"`
	Strength                float32 `json:"strength"`
}

type modelCompleted struct {
	Meta struct {
		NumImagesInLast24Hours int32 `json:"numImagesInLast24Hours"`
	} `json:"meta"`
	Images []struct {
		ImageKey string `json:"imageKey"`
		Prompt   string `json:"prompt"`
		Url      string `json:"url"`
		Loading  bool   `json:"loading"`
	} `json:"images"`
}

func Generation(ctx *gin.Context, req gpt.ChatGenerationRequest) {
	hash := sdio.SessionHash()
	var (
		cookie = ctx.GetString("token")
		domain = ctx.GetString("domain")

		modelType       = pkg.Config.GetString("playground.modelType")
		dreamBoothModel = pkg.Config.GetString("playground.dreamBoothModel")
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	var payload = modelPayload{
		BatchId:                 hash,
		CfgScale:                7,
		DreamBoothModel:         dreamBoothModel,
		Filter:                  dreamBoothModel,
		GenerateVariants:        false,
		GuidanceScale:           7,
		Width:                   1024,
		Height:                  1024,
		HighNoiseFrac:           0.8,
		InitImageFromPlayground: false,
		IsPrivate:               true,
		ModelType:               modelType,
		NegativePrompt:          "ugly, deformed, noisy, blurry, distorted, out of focus, bad anatomy, extra limbs, poorly drawn face, poorly drawn hands, missing fingers, ugly, deformed, noisy, blurry, distorted, out of focus, bad anatomy, extra limbs, poorly drawn face, poorly drawn hands, missing fingers, photo, realistic, text, watermark, signature, username, artist name",
		NumImages:               1,
		Prompt:                  req.Prompt,
		Sampler:                 9,
		Seed:                    int32(rand.Intn(100000000) + 329650152),
		StatusUUID:              uuid.NewString(),
		Steps:                   30,
		Strength:                1.45,
	}

	marshal, _ := json.Marshal(payload)
	response, err := fetch("", cookie, marshal)
	if err != nil {
		middle.ResponseWithE(ctx, -1, err)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		middle.ResponseWithE(ctx, -1, err)
		return
	}

	var mc modelCompleted
	if err = json.Unmarshal(data, &mc); err != nil {
		middle.ResponseWithE(ctx, -1, err)
		return
	}

	if len(mc.Images) == 0 {
		middle.ResponseWithV(ctx, -1, "generate images failed")
		return
	}

	file, err := saveImage(mc.Images[0].Url)
	if err != nil {
		middle.ResponseWithE(ctx, -1, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"created": time.Now().Unix(),
		"data": []map[string]string{
			{"url": fmt.Sprintf("%s/file/%s", domain, file)},
		},
	})
}

func saveImage(base64Encoding string) (file string, err error) {
	index := strings.Index(base64Encoding, ",")
	base64Encoding = base64Encoding[index+1:]
	decode, err := base64.StdEncoding.DecodeString(base64Encoding)
	if err != nil {
		return "", err
	}

	_, err = os.Stat("images")
	if os.IsNotExist(err) {
		_ = os.Mkdir("images", os.ModePerm)
	}

	tempFile, err := os.CreateTemp("images", "image-*.jpg")
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(decode)
	if err != nil {
		return "", err
	}

	file = tempFile.Name()
	return file, nil
}

func fetch(proxies, cookie string, marshal []byte) (*http.Response, error) {
	if !strings.Contains(cookie, "__Secure-next-auth.session-token=") {
		cookie = "__Secure-next-auth.session-token=" + cookie
	}

	client, err := common.NewHttpClient(proxies)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://playground.com/api/models", bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	h := request.Header
	h.Add("accept", "application/json")
	h.Add("content-type", "application/json")
	h.Add("origin", "https://playground.com")
	h.Add("referer", "https://playground.com/create")
	h.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.79")
	h.Add("cookie", cookie)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	return response, nil
}
