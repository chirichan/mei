package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"

// ===== 数据结构 =====

type LoginResp struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		User struct {
			Token string `json:"token"`
		} `json:"user"`
	} `json:"data"`
}

type SignResp struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// ===== 主函数 =====

func main() {
	userFlag := flag.String("user", "", "username")
	passFlag := flag.String("pass", "", "password")
	serveFlag := flag.Bool("serve", false, "run as daemon")
	timeFlag := flag.String("time", "09:00", "daily run time (HH:MM)")
	noFirst := flag.Bool("no-first", false, "skip first immediate run")

	flag.Parse()

	username := firstNonEmpty(*userFlag, os.Getenv("ZAI_USER"))
	password := firstNonEmpty(*passFlag, os.Getenv("ZAI_PASS"))

	if username == "" || password == "" {
		fmt.Println("Usage:")
		fmt.Println("  --user xxx --pass xxx [--serve --time 09:00]")
		fmt.Println("  or set env: ZAI_USER / ZAI_PASS")
		return
	}

	if *serveFlag {
		runDaemon(username, password, *timeFlag, !*noFirst)
	} else {
		runOnce(username, password)
	}
}

// ===== 执行模式 =====

func runOnce(username, password string) {
	fmt.Println("[INFO] 单次执行签到")
	if err := doSign(username, password); err != nil {
		fmt.Println("[ERROR]", err)
	}
}

func runDaemon(username, password, timeStr string, runImmediately bool) {
	fmt.Println("[INFO] daemon 模式启动")

	hour, min, err := parseTime(timeStr)
	if err != nil {
		fmt.Println("[ERROR] 时间格式错误，应为 HH:MM")
		return
	}

	// 优雅退出
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if runImmediately {
		fmt.Println("[INFO] 启动立即执行一次")
		doSign(username, password)
	}

	for {
		next := nextRunTime(hour, min)
		fmt.Println("[INFO] 下次执行时间:", next.Format(time.RFC3339))

		select {
		case <-ctx.Done():
			fmt.Println("[INFO] 收到退出信号，程序结束")
			return
		case <-time.After(time.Until(next)):
			fmt.Println("[INFO] 开始执行签到:", time.Now().Format(time.RFC3339))
			if err := doSign(username, password); err != nil {
				fmt.Println("[ERROR]", err)
			}
		}
	}
}

// ===== 时间计算 =====

// 解析 "09:00"
func parseTime(s string) (int, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time")
	}

	var h, m int
	_, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	if err != nil {
		return 0, 0, err
	}

	return h, m, nil
}

// 计算下一次执行时间
func nextRunTime(hour, min int) time.Time {
	now := time.Now()
	next := time.Date(
		now.Year(), now.Month(), now.Day(),
		hour, min, 0, 0,
		now.Location(),
	)

	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}

	return next
}

// ===== 核心业务 =====

func doSign(username, password string) error {
	token, err := login(username, password)
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	fmt.Println("[INFO] 登录成功")

	err = signIn(token)
	if err != nil {
		return fmt.Errorf("签到失败: %w", err)
	}

	return nil
}

func login(username, password string) (string, error) {
	baseURL := "https://i.zaimanhua.com/lpi/v1/login/passwd"

	params := url.Values{}
	params.Set("username", username)
	params.Set("passwd", password)

	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result LoginResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Errno != 0 {
		return "", fmt.Errorf("%s", result.Errmsg)
	}

	return result.Data.User.Token, nil
}

func signIn(token string) error {
	req, err := http.NewRequest(
		"POST",
		"https://i.zaimanhua.com/lpi/v1/task/sign_in",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result SignResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Errno == 0 {
		fmt.Println("[INFO] 签到成功")
		return nil
	}

	fmt.Println("[INFO]", result.Errmsg)
	return nil
}

// ===== 工具 =====

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
