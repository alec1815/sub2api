// +build ignore

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

	client := &http.Client{}

	// 1. 企业列表获取第一个企业 ID
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/admin/enterprises?page_size=1", nil)
	req.Header.Set("Authorization", token)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("List => HTTP %d: %s\n", resp.StatusCode, string(body[:min(300, len(body))]))

	// 2. 测试充值
	ts := time.Now().UnixNano()
	balanceBody := fmt.Sprintf(`{"balance":100,"operation":"add","notes":"test deposit %d"}`, ts%100000)
	req2, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/v1/admin/enterprises/1/balance", bytes.NewBufferString(balanceBody))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", token)
	resp2, _ := client.Do(req2)
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("Deposit(EID=1) => HTTP %d: %s\n", resp2.StatusCode, string(body2[:min(200, len(body2))]))

	// 3. 测试扣款
	withdrawBody := `{"balance":10,"operation":"subtract","notes":"test withdrawal"}`
	req3, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/v1/admin/enterprises/1/balance", bytes.NewBufferString(withdrawBody))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", token)
	resp3, _ := client.Do(req3)
	body3, _ := io.ReadAll(resp3.Body)
	resp3.Body.Close()
	fmt.Printf("Withdraw(EID=1) => HTTP %d: %s\n", resp3.StatusCode, string(body3[:min(200, len(body3))]))

	// 4. 测试 balance-history
	req4, _ := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/admin/enterprises/1/balance-history?page=1&page_size=5", nil)
	req4.Header.Set("Authorization", token)
	resp4, _ := client.Do(req4)
	body4, _ := io.ReadAll(resp4.Body)
	resp4.Body.Close()
	fmt.Printf("History(EID=1) => HTTP %d: %s\n", resp4.StatusCode, string(body4[:min(200, len(body4))]))
}

func min(a, b int) int { if a < b { return a }; return b }
