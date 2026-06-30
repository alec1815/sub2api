// +build ignore

// 企业创建接口测试工具
// 用法: cd backend && go run test_req.go
// 默认使用 127.0.0.1:8080，可通过环境变量修改

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		fmt.Println("请设置 TEST_TOKEN 环境变量")
		os.Exit(1)
	}
	token = "Bearer " + token

	host := os.Getenv("TEST_HOST")
	if host == "" {
		host = "http://127.0.0.1:8080"
	}

	client := &http.Client{}
	ts := time.Now().UnixNano()
	newEmail := fmt.Sprintf("test%d@local.dev", ts%100000)

	body := fmt.Sprintf(`{"name":"test-%d","contact_name":"c","contact_phone":"13800000001","contact_email":"c@qq.com","admin_email":"%s","admin_password":"admin123","admin_password_confirm":"admin123"}`, ts%1000, newEmail)

	req, _ := http.NewRequest("POST", host+"/api/v1/admin/enterprises", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	fmt.Printf("Request: %s\n", body)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("HTTP %d: %s\n", resp.StatusCode, string(respBody))
}
