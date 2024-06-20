package netx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClientI 定义 HTTP 客户端接口
type HTTPClientI interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string, headers map[string]string, result interface{}) error
	Post(url string, headers map[string]string, body, result interface{}) error
	Put(url string, headers map[string]string, body, result interface{}) error
	Delete(url string, headers map[string]string, body, result interface{}) error
}

// HTTPClient 实现 HTTP 客户端接口
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(timeout time.Duration) HTTPClientI {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do 发送 HTTP 请求
// req: 需要发送的 HTTP 请求
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// 通用的 HTTP 请求方法
// method: 请求方法（GET, POST, PUT, DELETE）
// url: 请求的 URL
// headers: 请求头
// body: 请求体，可以为 nil
// result: 返回结果，将响应体解析到该对象
func (c *HTTPClient) doRequest(method, url string, headers map[string]string, body, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("请求体序列化失败: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s 请求失败，状态码: %s", method, resp.Status)
	}

	if result != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取响应体失败: %v", err)
		}
		return json.Unmarshal(respBody, result)
	}

	return nil
}

// Get 发送 GET 请求
// url: 请求的 URL
// headers: 请求头
// result: 指针，返回结果，将响应体解析到该对象
// 返回错误信息
func (c *HTTPClient) Get(url string, headers map[string]string, result interface{}) error {
	return c.doRequest("GET", url, headers, nil, result)
}

// Post 发送 POST 请求
// url: 请求的 URL
// headers: 请求头
// body: 请求体
// result: 指针，返回结果，将响应体解析到该对象
func (c *HTTPClient) Post(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("POST", url, headers, body, result)
}

// Put 发送 PUT 请求
// url: 请求的 URL
// headers: 请求头
// body: 请求体
// result: 指针，返回结果，将响应体解析到该对象
func (c *HTTPClient) Put(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("PUT", url, headers, body, result)
}

// Delete 发送 DELETE 请求
// url: 请求的 URL
// headers: 请求头
// body: 请求体，可以为 nil
// result: 指针，返回结果，将响应体解析到该对象
func (c *HTTPClient) Delete(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("DELETE", url, headers, body, result)
}
