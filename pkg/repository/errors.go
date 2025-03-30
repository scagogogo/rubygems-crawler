package repository

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrInvalidRequest 请求参数无效
	ErrInvalidRequest = errors.New("invalid request parameters")

	// ErrNotFound 资源未找到
	ErrNotFound = errors.New("resource not found")

	// ErrServerError 服务器错误
	ErrServerError = errors.New("server error")

	// ErrRateLimited 请求被限流
	ErrRateLimited = errors.New("request rate limited")

	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")

	// ErrTimeout 请求超时
	ErrTimeout = errors.New("request timeout")

	// ErrNetworkFailure 网络故障
	ErrNetworkFailure = errors.New("network failure")
)

// APIError 表示API调用时遇到的错误
type APIError struct {
	// 错误原因
	Cause error

	// HTTP状态码
	StatusCode int

	// 请求URL
	URL string

	// 响应内容
	Response string
}

// 实现Error接口
func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status: %d, url: %s): %v", e.StatusCode, e.URL, e.Cause)
}

// 从HTTP响应创建APIError
func NewAPIError(resp *http.Response, body []byte, cause error) *APIError {
	return &APIError{
		Cause:      cause,
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
		Response:   string(body),
	}
}

// IsNotFound 检查错误是否为资源未找到
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return errors.Is(err, ErrNotFound)
}

// IsRateLimited 检查错误是否为请求被限流
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return errors.Is(err, ErrRateLimited)
}

// IsUnauthorized 检查错误是否为未授权
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return errors.Is(err, ErrUnauthorized)
}
