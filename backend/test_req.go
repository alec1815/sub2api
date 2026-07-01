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
	ts := time.Now().UnixNano()

	do := func(method, path, body string) (int, string) {
		var r io.Reader
		if body != "" {
			r = bytes.NewBufferString(body)
		}
		req, _ := http.NewRequest(method, "http://127.0.0.1:8080"+path, r)
		req.Header.Set("Authorization", token)
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := client.Do(req)
		if err != nil {
			return 0, err.Error()
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		s := string(b)
		if len(s) > 200 {
			s = s[:200] + "..."
		}
		return resp.StatusCode, s
	}

	fmt.Println("=== 1. 列表 ===")
	code, body := do("GET", "/api/v1/admin/enterprises?page_size=1", "")
	fmt.Printf("List => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 2. 充值 ===")
	code, body = do("POST", "/api/v1/admin/enterprises/1/balance", fmt.Sprintf(`{"balance":50,"operation":"add","notes":"dep-%d"}`, ts%1000))
	fmt.Printf("Deposit => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 3. 扣款 ===")
	code, body = do("POST", "/api/v1/admin/enterprises/1/balance", `{"balance":5,"operation":"subtract","notes":"withdrawal"}`)
	fmt.Printf("Withdraw => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 4. 充值记录 ===")
	code, body = do("GET", "/api/v1/admin/enterprises/1/balance-history?page_size=2", "")
	fmt.Printf("History => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 5. API密钥 ===")
	code, body = do("GET", "/api/v1/admin/enterprises/1/api-keys?page_size=2", "")
	fmt.Printf("APIKeys => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 6. 平台限额 (GET) ===")
	code, body = do("GET", "/api/v1/admin/enterprises/1/platform-quotas", "")
	fmt.Printf("Quotas(GET) => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 7. 平台限额 (PUT) ===")
	code, body = do("PUT", "/api/v1/admin/enterprises/1/platform-quotas", `{"quotas":[{"platform":"openai","daily_limit_usd":100}]}`)
	fmt.Printf("Quotas(PUT) => HTTP %d: %s\n", code, body)

	fmt.Println("\n=== 8. 验证限额已设置 ===")
	code, body = do("GET", "/api/v1/admin/enterprises/1/platform-quotas", "")
	fmt.Printf("Quotas(GET2) => HTTP %d: %s\n", code, body)
}
