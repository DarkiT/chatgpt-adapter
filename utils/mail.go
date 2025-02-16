package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

const (
	maxPollAttempts = 120
	pollInterval    = 2 * time.Second
	baseURL         = "https://etempmail.com"
	cursorBaseURL   = "https://www.cursor.com"
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

// UserInfo 定义了用户信息的结构体，与 JSON 数据字段对应
type UserInfo struct {
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Name          string    `json:"name"`
	Sub           string    `json:"sub"`
	UpdatedAt     time.Time `json:"updated_at"`
	Picture       *string   `json:"picture"` // 当 JSON 中为 null 时，该字段将为 nil
}

// ModelUsage 定义了各模型的使用情况
type ModelUsage struct {
	NumRequests      int `json:"numRequests"`
	NumRequestsTotal int `json:"numRequestsTotal"`
	NumTokens        int `json:"numTokens"`
	// 指针类型，用于处理 JSON 中的 null 值
	MaxRequestUsage *int `json:"maxRequestUsage"`
	MaxTokenUsage   *int `json:"maxTokenUsage"`
}

// GPTUsageData 定义了所有模型的使用数据以及统计起始时间
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

func (c *EmailClient) GetEmailAddress() (*EmailAddress, error) {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/getEmailAddress", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
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
	req, err := http.NewRequest(http.MethodPost, baseURL+"/getInbox", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.AddCookie(c.ciSession)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
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

func (c *EmailClient) GetUsage(session string) (*GPTUsageData, error) {
	req, err := http.NewRequest(http.MethodGet, cursorBaseURL+"/api/auth/me", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Cookie", "WorkosCursorSessionToken="+session)
	req.Header.Add("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var userInfo UserInfo

	jsonErr := json.Unmarshal(body, &userInfo)
	if jsonErr != nil {
		return nil, fmt.Errorf("解析响应失败: %w", jsonErr)
	}

	req, err = http.NewRequest(http.MethodGet, cursorBaseURL+"/api/usage?user="+userInfo.Sub, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Cookie", "WorkosCursorSessionToken="+session)
	req.Header.Add("User-Agent", userAgent)

	resp1, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp1.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var usage GPTUsageData
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &usage, nil
}

// extractVerificationCode 从传入的 HTML 字符串中查找并返回 6 位验证码
func extractVerificationCode(htmlContent string) (string, error) {
	var code string
	// 构建正则表达式，匹配 "Your one-time code is" 后跟随的数字
	re := regexp.MustCompile(`Your one-time code is\s+(\d+)\.`)
	// 在 HTML 内容中查找匹配项
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		code = matches[1]
	}

	if code == "" {
		return "", fmt.Errorf("未找到验证码")
	}
	return code, nil
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

/*
func main() {

	client := NewEmailClient()

	usage, err := client.GetUsage(session)
	if err != nil {
		return
	}

	slog.Infof("GPT-4 使用数据Num：%d Token: %d", usage.GPT4.NumRequests, usage.GPT4.NumTokens)
	slog.Infof("GPT-3.5-Turbo 使用数据：%d", usage.GPT35Turbo.NumRequests)
	slog.Infof("GPT-4-32k 使用数据：%d", usage.GPT4_32k.NumRequests)
	slog.Infof("统计起始时间：%s", usage.StartOfMonth.Format(time.DateTime))

	// 获取邮箱地址
	emailAddr, err := client.GetEmailAddress()
	if err != nil {
		slog.Fatalf("初始化失败: %v", err)
	}
	slog.Infof("成功获取邮箱地址: %s", emailAddr.Address)

	// 轮询收件箱
	var emails []Email
	for attempt := 1; attempt <= maxPollAttempts; attempt++ {
		emails, err = client.PollInbox()
		if err != nil {
			slog.Errorf("轮询失败（尝试次数 %d）: %v", attempt, err)
			continue
		}

		if len(emails) > 0 {
			fmt.Printf("第 %d 次尝试成功收到邮件\n", attempt)
			break
		}

		slog.Infof("第 %d 次轮询: 收件箱为空", attempt)
		time.Sleep(pollInterval)
	}

	if len(emails) == 0 {
		slog.Fatal("超过最大轮询次数仍未收到邮件")
	}

	// 输出邮件内容
	for i, email := range emails {
		code, err := extractVerificationCode(email.Body)
		if err != nil {
			code = "未找到验证码"
		}
		slog.Infof("\n邮件 #%d:\n主题: %s\n发件人: %s\n日期: %s\n验证码: %s", i+1, email.Subject, email.From, email.Date, code)
	}
}
*/
