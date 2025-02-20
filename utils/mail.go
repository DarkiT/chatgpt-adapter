package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/darkit/slog"
)

const (
	maxPollAttempts = 60
	pollInterval    = 2 * time.Second
	baseURL         = "https://etempmail.com"
	apiBaseURL      = "https://www.cursor.com"
	tokenApi        = "https://cursor.ccopilot.org/api/get_next_token.php"
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML,like Gecko) Chrome/108.0.0.0 Safari/537.36"
)

type EmailAddress struct {
	ID           string `json:"id"`
	Address      string `json:"address"`
	CreationTime string `json:"creation_time"`
	RecoverKey   string `json:"recover_key"`
}

type Email struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	Date    string `json:"date"`
	Body    string `json:"body"`
}

type EmailClient struct {
	httpClient *http.Client
	ciSession  *http.Cookie
}

type UserInfo struct {
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Name          string    `json:"name"`
	Sub           string    `json:"sub"`
	UpdatedAt     time.Time `json:"updated_at"`
	Picture       *string   `json:"picture"`
}

type ModelUsage struct {
	NumRequests      int  `json:"numRequests"`
	NumRequestsTotal int  `json:"numRequestsTotal"`
	NumTokens        int  `json:"numTokens"`
	MaxRequestUsage  *int `json:"maxRequestUsage"`
	MaxTokenUsage    *int `json:"maxTokenUsage"`
}

type LoginDeepControlRequest struct {
	UUID      string `json:"uuid"`
	Challenge string `json:"challenge"`
}

type GPTUsageData struct {
	GPT4         ModelUsage `json:"gpt-4"`
	GPT35Turbo   ModelUsage `json:"gpt-3.5-turbo"`
	GPT4_32k     ModelUsage `json:"gpt-4-32k"`
	StartOfMonth time.Time  `json:"startOfMonth"`
}

func NewEmailClient() *EmailClient {
	return &EmailClient{
		httpClient: &http.Client{},
	}
}

func (c *EmailClient) doRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return c.httpClient.Do(req)
}

func (c *EmailClient) GetEmailAddress() (*EmailAddress, error) {
	resp, err := c.doRequest(http.MethodPost, baseURL+"/getEmailAddress", nil, map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	var emailAddr EmailAddress
	if err := json.NewDecoder(resp.Body).Decode(&emailAddr); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	c.ciSession = findCookie(resp.Cookies(), "ci_session")
	if c.ciSession == nil {
		return nil, fmt.Errorf("未找到ci_session cookie")
	}

	return &emailAddr, nil
}

func (c *EmailClient) PollInbox() ([]Email, error) {
	resp, err := c.doRequest(http.MethodPost, baseURL+"/getInbox", nil, map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Cookie":       c.ciSession.String(),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	var emails []Email
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return emails, nil
}

func (c *EmailClient) GetUsage(userID, cookie string) (*GPTUsageData, error) {
	headers := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Cookie":       "WorkosCursorSessionToken=" + fmt.Sprintf("%s%%3A%%3A%s", userID, cookie),
		"User-Agent":   userAgent,
	}

	resp, err := c.doRequest(http.MethodGet, apiBaseURL+"/api/usage?user="+userID, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	var usage GPTUsageData
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &usage, nil
}

func (c *EmailClient) GetVerificationCode() (string, error) {
	for attempt := 1; attempt <= maxPollAttempts; attempt++ {
		emails, err := c.PollInbox()
		if err != nil {
			slog.Warn("邮箱轮询失败", "尝试次数", attempt, "错误", err.Error())
			time.Sleep(pollInterval)
			continue
		}

		for _, email := range emails {
			if code, err := extractVerificationCode(email.Body); err == nil {
				slog.Info("成功提取验证码", "尝试次数", attempt, "发件人", email.From, "主题", email.Subject)
				return code, nil
			}
		}

		slog.Debug("未找到有效验证码", "当前邮件数", len(emails), "剩余尝试", maxPollAttempts-attempt)
		time.Sleep(pollInterval)
	}

	return "", fmt.Errorf("经%d次尝试仍未获取到验证码", maxPollAttempts)
}

func (c *EmailClient) LoginDeepControl(inputURL string) error {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return err
	}

	queryParams := parsedURL.Query()
	body := LoginDeepControlRequest{
		UUID:      queryParams.Get("uuid"),
		Challenge: queryParams.Get("challenge"),
	}

	if body.UUID == "" || body.Challenge == "" {
		return fmt.Errorf("非法请求")
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Referer":      fmt.Sprintf("/cn/loginDeepControl?challenge=%s&uuid=%s&mode=login", body.Challenge, body.UUID),
		"Cookie":       "WorkosCursorSessionToken=" + GetAuthToken().Cookie,
		"User-Agent":   userAgent,
	}

	resp, err := c.doRequest("POST", apiBaseURL+"/api/auth/loginDeepCallbackControl", bytes.NewBuffer(jsonData), headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func (c *EmailClient) SyncSessionToken(token string) error {
	type tokenData struct {
		Code string `json:"code"`
		Data string `json:"data"`
	}

	if token == "" {
		return fmt.Errorf("token不能为空")
	}
	slog.Info("正在获新的SessionToken...")
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   userAgent,
	}
	data := map[string]string{
		"accessCode":    token,
		"cursorVersion": "0.45.11",
		"scriptVersion": "2025020801",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := c.doRequest("POST", tokenApi, bytes.NewBuffer(jsonData), headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if str == nil {
		var tmp tokenData
		if err = json.Unmarshal(str, &tmp); err != nil {
			return err
		}
		config := &ConfigRequest{
			Cookie: tmp.Data,
		}
		sum, ok := GenChecksum(config)
		if ok {
			config.Checksum = sum
		}
		if ok, _ = AuthToken(config, true); !ok {
			return fmt.Errorf("自动更新Token失败")
		}
		slog.Info("更新SessionToken成功")
		return nil
	}
	return err
}

func extractVerificationCode(htmlContent string) (string, error) {
	re := regexp.MustCompile(`Your one-time code is\s+(\d+)\.`)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("未找到验证码")
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
