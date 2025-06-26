package netx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func newRetryableClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.HTTPClient.Timeout = 3 * time.Second
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 5 * time.Second
	client.Logger = nil // 禁用日志输
	return client
}

func getReq(url string) (string, error) {
	client := newRetryableClient()
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func postReq(url string, body []byte) (string, error) {
	client := newRetryableClient()
	req, err := retryablehttp.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("POST 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}

func Get() {
	body, err := getReq("https://httpbin.org/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET 响应内容:", body)
}

type PostBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Post() {
	postBody := PostBody{
		Name: "hogwarts",
		Age:  18,
	}
	jsonBody, err := json.Marshal(postBody)
	if err != nil {
		log.Fatal("序列化请求体失败:", err)
	}
	body, err := postReq("https://httpbin.org/post", jsonBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST 响应内容:", body)
}
