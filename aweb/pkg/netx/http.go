package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sony/gobreaker"
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

func newBreaker() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "my-circuit-breaker",
		MaxRequests: 3,                // 在半开（Half-Open）状态下，允许通过的最大请求数。
		Interval:    60 * time.Second, // 统计失败率的时间窗口，超过该时间会重置统计计数。
		Timeout:     10 * time.Second, // 断路器从打开（Open）状态到半开（Half-Open）状态的等待时间。
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 失败率超过60%则跳闸
			return counts.ConsecutiveFailures > 5 ||
				(counts.Requests >= 10 && float64(counts.TotalFailures)/float64(counts.Requests) > 0.6)
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("断路器状态变化: %s -> %s\n", from.String(), to.String())
		},
	}
	cb := gobreaker.NewCircuitBreaker(settings)
	return cb
}

type GRestClient struct {
	cb     *gobreaker.CircuitBreaker
	client *retryablehttp.Client
}

func NewGRestClient() *GRestClient {
	return &GRestClient{
		cb:     newBreaker(),
		client: newRetryableClient(),
	}
}

func (c *GRestClient) getReq(url string) (string, error) {
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	respRaw, err := c.cb.Execute(func() (interface{}, error) {
		return c.client.Do(req)
	})
	if err != nil {
		return "", fmt.Errorf("GET 请求失败: %w", err)
	}
	resp, ok := respRaw.(*http.Response)
	if !ok {
		return "", fmt.Errorf("响应类型断言失败")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (c *GRestClient) postReq(url string, body []byte) (string, error) {
	req, err := retryablehttp.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	respRaw, err := c.cb.Execute(func() (interface{}, error) {
		return c.client.Do(req)
	})
	if err != nil {
		return "", fmt.Errorf("POST 请求失败: %w", err)
	}
	resp, ok := respRaw.(*http.Response)
	if !ok {
		return "", fmt.Errorf("响应类型断言失败")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}

func GetTest(GRestClient *GRestClient) {
	body, err := GRestClient.getReq("https://httpbin.org/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET 响应内容:", body)
}

type PostBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func PostTest(GRestClient *GRestClient) {
	postBody := PostBody{
		Name: "hogwarts",
		Age:  18,
	}
	jsonBody, err := json.Marshal(postBody)
	if err != nil {
		log.Fatal("序列化请求体失败:", err)
	}
	body, err := GRestClient.postReq("https://httpbin.org/post", jsonBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST 响应内容:", body)
}
