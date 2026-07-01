// +build ignore
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("TEST_TOKEN")
	if token == "" { fmt.Println("set TEST_TOKEN"); os.Exit(1) }
	token = "Bearer " + token
	do := func(method, path string) (int, string) {
		req, _ := http.NewRequest(method, "http://127.0.0.1:8080"+path, nil)
		req.Header.Set("Authorization", token)
		resp, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		s := string(body); if len(s) > 400 { s = s[:400]+"..." }
		return resp.StatusCode, s
	}
	fmt.Println("=== /enterprise/usage/models ===")
	c, b := do("GET", "/api/v1/enterprise/usage/models?start_date=2026-06-30&end_date=2026-07-01")
	fmt.Printf("HTTP %d: %s\n", c, b)
	fmt.Println("\n=== /enterprise/usage/trend ===")
	c, b = do("GET", "/api/v1/enterprise/usage/trend?start_date=2026-06-30&end_date=2026-07-01&granularity=day&limit=12")
	fmt.Printf("HTTP %d: %s\n", c, b)
	fmt.Println("\n=== /enterprise/keys ===")
	c, b = do("GET", "/api/v1/enterprise/keys?page_size=2")
	fmt.Printf("HTTP %d: %s\n", c, b)
	fmt.Println("\n=== /enterprise/members ===")
	c, b = do("GET", "/api/v1/enterprise/members?page_size=2")
	fmt.Printf("HTTP %d: %s\n", c, b)
}
