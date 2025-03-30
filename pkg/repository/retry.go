package repository

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/crawler-go-go-go/go-requests"
)

const (
	// DefaultRetryAttempts 默认重试次数
	DefaultRetryAttempts = 3

	// DefaultRetryWaitTime 默认重试等待时间
	DefaultRetryWaitTime = 1 * time.Second

	// DefaultRetryMaxWaitTime 默认最大重试等待时间
	DefaultRetryMaxWaitTime = 30 * time.Second
)

// RetryOptions 重试选项
type RetryOptions struct {
	// 重试次数
	MaxAttempts int

	// 初始等待时间
	WaitTime time.Duration

	// 最大等待时间
	MaxWaitTime time.Duration

	// 是否使用指数退避算法
	UseExponentialBackoff bool

	// 自定义重试条件
	ShouldRetry func(resp *http.Response, err error) bool
}

// NewDefaultRetryOptions 创建默认重试选项
func NewDefaultRetryOptions() *RetryOptions {
	return &RetryOptions{
		MaxAttempts:           DefaultRetryAttempts,
		WaitTime:              DefaultRetryWaitTime,
		MaxWaitTime:           DefaultRetryMaxWaitTime,
		UseExponentialBackoff: true,
		ShouldRetry: func(resp *http.Response, err error) bool {
			// 如果有错误，总是重试
			if err != nil {
				return true
			}

			// 对于特定的HTTP状态码，进行重试
			if resp != nil {
				switch resp.StatusCode {
				case http.StatusTooManyRequests, // 429
					http.StatusInternalServerError, // 500
					http.StatusBadGateway,          // 502
					http.StatusServiceUnavailable,  // 503
					http.StatusGatewayTimeout:      // 504
					return true
				}
			}

			return false
		},
	}
}

// WithMaxAttempts 设置最大重试次数
func (o *RetryOptions) WithMaxAttempts(attempts int) *RetryOptions {
	o.MaxAttempts = attempts
	return o
}

// WithWaitTime 设置初始等待时间
func (o *RetryOptions) WithWaitTime(waitTime time.Duration) *RetryOptions {
	o.WaitTime = waitTime
	return o
}

// WithMaxWaitTime 设置最大等待时间
func (o *RetryOptions) WithMaxWaitTime(maxWaitTime time.Duration) *RetryOptions {
	o.MaxWaitTime = maxWaitTime
	return o
}

// WithExponentialBackoff 设置是否使用指数退避算法
func (o *RetryOptions) WithExponentialBackoff(use bool) *RetryOptions {
	o.UseExponentialBackoff = use
	return o
}

// WithShouldRetry 设置自定义重试条件
func (o *RetryOptions) WithShouldRetry(shouldRetry func(resp *http.Response, err error) bool) *RetryOptions {
	o.ShouldRetry = shouldRetry
	return o
}

// SendRequestWithRetry 发送带重试功能的请求
func SendRequestWithRetry[Request any, Response any](
	ctx context.Context,
	options *requests.Options[Request, Response],
	retryOptions *RetryOptions,
) (Response, error) {
	var lastErr error
	var lastResp Response

	// 如果未提供重试选项，使用默认值
	if retryOptions == nil {
		retryOptions = NewDefaultRetryOptions()
	}

	for attempt := 0; attempt < retryOptions.MaxAttempts; attempt++ {
		// 如果不是第一次尝试，等待一段时间
		if attempt > 0 {
			waitTime := retryOptions.WaitTime

			// 如果使用指数退避，则指数增加等待时间
			if retryOptions.UseExponentialBackoff {
				factor := 1 << uint(attempt-1)
				waitTime = time.Duration(float64(waitTime) * float64(factor))
				if waitTime > retryOptions.MaxWaitTime {
					waitTime = retryOptions.MaxWaitTime
				}
			}

			// 等待一段时间后重试
			select {
			case <-time.After(waitTime):
				// 继续执行
			case <-ctx.Done():
				// 上下文被取消，停止重试
				var zero Response
				return zero, ctx.Err()
			}
		}

		// 执行请求
		resp, err := requests.SendRequest[Request, Response](ctx, options)

		// 检查是否需要重试
		shouldRetry := false
		if err != nil {
			lastErr = err
			shouldRetry = true
		} else {
			// 请求成功，返回结果
			return resp, nil
		}

		// 如果不需要重试，直接返回结果
		if !shouldRetry {
			return resp, nil
		}

		// 记录最后一次响应
		lastResp = resp
	}

	// 达到最大重试次数，返回最后一次的错误
	if lastErr != nil {
		return lastResp, errors.New("max retry attempts reached: " + lastErr.Error())
	}

	return lastResp, nil
}
