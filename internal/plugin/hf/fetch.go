package hf

import (
	"chatgpt-adapter/internal/common"
	"chatgpt-adapter/internal/plugin"
	"chatgpt-adapter/logger"
	"chatgpt-adapter/pkg"
	"fmt"
	"github.com/bincooo/emit.io"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const (
	negative  = "(deformed eyes, nose, ears, nose, leg, head), bad anatomy, ugly"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0"
)

func Ox000(ctx *gin.Context, model, samples, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		baseUrl = "https://prodia-fast-stable-diffusion.hf.space"
		domain  = pkg.Config.GetString("domain")
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	response, err := emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		POST(baseUrl+"/queue/join").
		JHeader().
		Header("User-Agent", userAgent).
		Body(map[string]interface{}{
			"data": []interface{}{
				message + ", {{{{by famous artist}}}, beautiful, 4k",
				negative,
				model,
				25,
				samples,
				10,
				1024,
				1024,
				-1,
			},
			"fn_index":     0,
			"trigger_id":   15,
			"session_hash": hash,
		}).
		DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return
	}

	logger.Info(emit.TextResponse(response))
	response, err = emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		GET(baseUrl+"/queue/data").
		Query("session_hash", hash).
		Header("User-Agent", userAgent).
		DoS(http.StatusOK)
	if err != nil {
		return
	}

	c, err := emit.NewGio(common.GetGinContext(ctx), response)
	if err != nil {
		return
	}

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		d := j.Output.Data
		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}
		result := d[0].(map[string]interface{})
		value = result["url"].(string)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}

func Ox001(ctx *gin.Context, model, samples, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		baseUrl = "wss://prodia-sdxl-stable-diffusion-xl.hf.space"
		domain  = pkg.Config.GetString("domain")
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	conn, response, err := emit.SocketBuilder(plugin.HTTPClient).
		Proxies(proxies).
		URL(baseUrl + "/queue/join").
		DoS(http.StatusSwitchingProtocols)
	if err != nil {
		return
	}
	defer response.Body.Close()

	c, err := emit.NewGio(common.GetGinContext(ctx), conn)
	if err != nil {
		return
	}

	c.Event("send_hash", func(j emit.JoinEvent) interface{} {
		return map[string]interface{}{
			"fn_index":     0,
			"session_hash": hash,
		}
	})

	c.Event("send_data", func(j emit.JoinEvent) interface{} {
		return map[string]interface{}{
			"fn_index":     0,
			"session_hash": hash,
			"data": []interface{}{
				message + ", {{{{by famous artist}}}, beautiful, 4k",
				negative,
				model,
				25,
				samples,
				10,
				1024,
				1024,
				-1,
			},
		}
	})

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		var file string
		d := j.Output.Data

		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		if file, err = common.SaveBase64(d[0].(string), "png"); err != nil {
			c.Failed(fmt.Errorf("image save failed: %s", j.InitialBytes))
			return
		}

		value = fmt.Sprintf("%s/file/%s", domain, file)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}

func Ox002(ctx *gin.Context, model, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		baseUrl = "https://mukaist-dalle-4k.hf.space"
	)

	if u := pkg.Config.GetString("hf.dalle-4k"); u != "" {
		baseUrl = u
	}

	response, err := emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		POST(baseUrl+"/queue/join").
		JHeader().
		Body(map[string]interface{}{
			"data": []interface{}{
				message,
				negative,
				true,
				model,
				30,
				1024,
				1024,
				6,
				true,
			},
			"fn_index":     3,
			"trigger_id":   6,
			"session_hash": hash,
		}).
		DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return
	}

	logger.Info(emit.TextResponse(response))
	response, err = emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		GET(baseUrl+"/queue/data").
		Query("session_hash", hash).
		DoC(emit.Status(http.StatusOK), emit.IsSTREAM)
	if err != nil {
		return
	}

	c, err := emit.NewGio(common.GetGinContext(ctx), response)
	if err != nil {
		return
	}

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		d := j.Output.Data

		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		values, ok := d[0].([]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}
		if len(values) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		v, ok := values[0].(map[string]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		value = v["image"].(map[string]interface{})["url"].(string)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}

func Ox003(ctx *gin.Context, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		domain  = pkg.Config.GetString("domain")
		baseUrl = "https://ehristoforu-dalle-3-xl-lora-v2.hf.space"
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	if u := pkg.Config.GetString("hf.dalle-3-xl"); u != "" {
		baseUrl = u
	}

	response, err := emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		POST(baseUrl+"/queue/join").
		Header("Origin", baseUrl).
		Header("Referer", baseUrl+"/?__theme=light").
		Header("User-Agent", userAgent).
		Header("Accept-Language", "en-US,en;q=0.9").
		JHeader().
		Body(map[string]interface{}{
			"data": []interface{}{
				message + ", {{{{by famous artist}}}, beautiful, 4k",
				negative + ", extra limb, missing limb, floating limbs, (mutated hands and fingers:1.4), disconnected limbs, mutation, mutated, ugly, disgusting, blurry, amputation",
				true,
				r.Intn(51206501) + 1100000000,
				1024,
				1024,
				12,
				true,
			},
			"fn_index":     3,
			"trigger_id":   6,
			"session_hash": hash,
		}).DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return "", err
	}
	logger.Info(emit.TextResponse(response))

	response, err = emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		GET(baseUrl+"/queue/data").
		Query("session_hash", hash).
		Header("Origin", baseUrl).
		Header("Referer", baseUrl+"/?__theme=light").
		Header("User-Agent", userAgent).
		Header("Accept", "text/event-stream").
		Header("Accept-Language", "en-US,en;q=0.9").
		DoC(emit.Status(http.StatusOK), emit.IsSTREAM)
	if err != nil {
		return "", err
	}

	c, err := emit.NewGio(ctx.Request.Context(), response)
	if err != nil {
		return "", err
	}

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		d := j.Output.Data

		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		values, ok := d[0].([]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		if len(values) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		v, ok := values[0].(map[string]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		info, ok := v["image"].(map[string]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		// 锁环境了，只能先下载下来
		value, err = common.Download(plugin.HTTPClient, proxies, info["url"].(string), "png", nil, map[string]string{
			"User-Agent":      userAgent,
			"Accept-Language": "en-US,en;q=0.9",
			"Origin":          baseUrl,
			"Referer":         baseUrl + "/?__theme=light",
		})
		if err != nil {
			c.Failed(fmt.Errorf("image download failed: %v", err))
			return
		}

		value = fmt.Sprintf("%s/file/%s", domain, value)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}

// 潦草漫画的风格
func Ox004(ctx *gin.Context, model, samples, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		domain  = pkg.Config.GetString("domain")
		baseUrl = "https://cagliostrolab-animagine-xl-3-1.hf.space"
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	if u := pkg.Config.GetString("hf.animagine-xl-3.1"); u != "" {
		baseUrl = u
	}
	jar, _ := emit.NewCookieJar(baseUrl, "")
	response, err := emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		POST(baseUrl+"/run/predict").
		CookieJar(jar).
		Header("Origin", baseUrl).
		Header("Referer", baseUrl+"/?__theme=light").
		Header("User-Agent", userAgent).
		Header("Accept-Language", "en-US,en;q=0.9").
		JHeader().
		Body(map[string]interface{}{
			"data":         []interface{}{0, true},
			"event_data":   nil,
			"fn_index":     4,
			"trigger_id":   49,
			"session_hash": hash,
		}).DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return "", err
	}
	logger.Info(emit.TextResponse(response))

	response, err = emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		POST(baseUrl+"/queue/join").
		CookieJar(jar).
		Header("Origin", baseUrl).
		Header("Referer", baseUrl+"/?__theme=light").
		Header("User-Agent", userAgent).
		Header("Accept-Language", "en-US,en;q=0.9").
		JHeader().
		Body(map[string]interface{}{
			"data": []interface{}{
				message + ", {{{{by famous artist}}}, beautiful, 4k",
				negative,
				r.Intn(9068457) + 300000000,
				1024,
				1024,
				8,
				35,
				samples,
				"1024 x 1024",
				model,
				"Light v3.1",
				true,
				0.55,
				1.5,
				false,
			},
			"fn_index":     5,
			"trigger_id":   49,
			"session_hash": hash,
		}).DoC(emit.Status(http.StatusOK), emit.IsJSON)
	if err != nil {
		return "", err
	}
	logger.Info(emit.TextResponse(response))

	response, err = emit.ClientBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		GET(baseUrl+"/queue/data").
		CookieJar(jar).
		Query("session_hash", hash).
		Header("Origin", baseUrl).
		Header("Referer", baseUrl+"/?__theme=light").
		Header("User-Agent", userAgent).
		Header("Accept", "text/event-stream").
		Header("Accept-Language", "en-US,en;q=0.9").
		DoC(emit.Status(http.StatusOK), emit.IsSTREAM)
	if err != nil {
		return "", err
	}

	c, err := emit.NewGio(common.GetGinContext(ctx), response)
	if err != nil {
		return "", err
	}

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		d := j.Output.Data

		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		values, ok := d[0].([]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		if len(values) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		v, ok := values[0].(map[string]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		info, ok := v["image"].(map[string]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		// 锁环境了，只能先下载下来
		value, err = common.Download(plugin.HTTPClient, proxies, info["url"].(string), "png", jar, map[string]string{
			"User-Agent":      userAgent,
			"Accept-Language": "en-US,en;q=0.9",
			"Origin":          baseUrl,
			"Referer":         baseUrl + "/?__theme=light",
		})
		if err != nil {
			c.Failed(fmt.Errorf("image download failed: %v", err))
			return
		}

		value = fmt.Sprintf("%s/file/%s", domain, value)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}

func google(ctx *gin.Context, model, message string) (value string, err error) {
	var (
		hash    = emit.GioHash()
		proxies = ctx.GetString("proxies")
		baseUrl = "wss://google-sdxl.hf.space"
		domain  = pkg.Config.GetString("domain")
	)

	if domain == "" {
		domain = fmt.Sprintf("http://127.0.0.1:%d", ctx.GetInt("port"))
	}

	conn, response, err := emit.SocketBuilder(plugin.HTTPClient).
		Proxies(proxies).
		Context(common.GetGinContext(ctx)).
		URL(baseUrl + "/queue/join").
		DoS(http.StatusSwitchingProtocols)
	if err != nil {
		return
	}
	defer response.Body.Close()

	c, err := emit.NewGio(common.GetGinContext(ctx), conn)
	if err != nil {
		return
	}

	c.Event("send_hash", func(j emit.JoinEvent) interface{} {
		return map[string]interface{}{
			"fn_index":     2,
			"session_hash": hash,
		}
	})

	c.Event("send_data", func(j emit.JoinEvent) interface{} {
		return map[string]interface{}{
			"fn_index":     2,
			"session_hash": hash,
			"data": []interface{}{
				message + ", {{{{by famous artist}}}, beautiful, 4k",
				negative,
				25,
				model,
			},
		}
	})

	c.Event("process_completed", func(j emit.JoinEvent) (_ interface{}) {
		var file string
		d := j.Output.Data

		if len(d) == 0 {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		values, ok := d[0].([]interface{})
		if !ok {
			c.Failed(fmt.Errorf("image generate failed: %s", j.InitialBytes))
			return
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if file, err = common.SaveBase64(values[r.Intn(len(values))].(string), "jpg"); err != nil {
			c.Failed(fmt.Errorf("image save failed: %s", j.InitialBytes))
			return
		}

		value = fmt.Sprintf("%s/file/%s", domain, file)
		return
	})

	err = c.Do()
	if err == nil && value == "" {
		err = fmt.Errorf("image generate failed")
	}
	return
}
