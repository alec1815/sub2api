// +build ignore
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		fmt.Println("请设置 TEST_TOKEN 环境变量 (需要企业管理员 token)")
		os.Exit(1)
	}
	token = "Bearer " + token

	// 1. 测试获取用户信息
	do := func(path string) (int, string) {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8080"+path, nil)
		req.Header.Set("Authorization", token)
		resp, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp.StatusCode, string(body)
	}

	// 获取用户信息
	fmt.Println("=== 1. 用户信息 ===")
	code, body := do("/api/v1/user/me")
	fmt.Printf("GET /user/me => %d\n", code)
	var me struct {
		Data struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}
	json.Unmarshal([]byte(body), &me)
	fmt.Printf("  ID=%d Email=%s Role=%s\n", me.Data.ID, me.Data.Email, me.Data.Role)

	// 2. 测试 GET /enterprise/profile
	fmt.Println("\n=== 2. 企业 Profile ===")
	code, body = do("/api/v1/enterprise/profile")
	fmt.Printf("GET /enterprise/profile => %d\n", code)
	s := string(body)
	if len(s) > 200 { s = s[:200] + "..." }
	fmt.Printf("  %s\n", s)

	// 3. 测试企业密钥接口（需要企业管理员权限）
	fmt.Println("\n=== 3. 企业密钥 ===")
	code, body = do("/api/v1/enterprise/keys")
	fmt.Printf("GET /enterprise/keys => %d\n", code)
	s = string(body)
	if len(s) > 150 { s = s[:150] + "..." }
	fmt.Printf("  %s\n", s)

	// 4. 解码 JWT token 看 claims
	fmt.Println("\n=== 4. JWT Claims ===")
	parts := strings.Split(os.Getenv("TEST_TOKEN"), ".")
	if len(parts) == 3 {
		fmt.Printf("Payload (base64): %s...\n", parts[1][:min(50, len(parts[1]))])
	}
}

func min(a, b int) int { if a < b { return a }; return b }
