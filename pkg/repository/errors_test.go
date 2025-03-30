package repository

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试API错误的创建
func TestNewAPIError(t *testing.T) {
	// 创建一个测试请求和响应
	req, _ := http.NewRequest("GET", "https://example.com/test", nil)
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Request:    req,
	}

	// 创建一个API错误
	cause := errors.New("test error")
	body := []byte("Not found")
	apiErr := NewAPIError(resp, body, cause)

	// 验证错误字段
	assert.Equal(t, cause, apiErr.Cause, "原始错误应该被保存")
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode, "状态码应该被保存")
	assert.Equal(t, "https://example.com/test", apiErr.URL, "URL应该被保存")
	assert.Equal(t, "Not found", apiErr.Response, "响应体应该被保存")
}

// 测试错误字符串表示
func TestAPIError_Error(t *testing.T) {
	apiErr := &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusInternalServerError,
		URL:        "https://example.com/test",
		Response:   "Server error",
	}

	errorStr := apiErr.Error()
	assert.Contains(t, errorStr, "500", "错误字符串应包含状态码")
	assert.Contains(t, errorStr, "https://example.com/test", "错误字符串应包含URL")
	assert.Contains(t, errorStr, "test error", "错误字符串应包含原始错误信息")
}

// 测试NotFound错误判断
func TestIsNotFound(t *testing.T) {
	// 测试直接的ErrNotFound
	assert.True(t, IsNotFound(ErrNotFound), "ErrNotFound应该被识别为NotFound")

	// 测试404状态码的API错误
	apiErr := &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusNotFound,
		URL:        "https://example.com/test",
	}
	assert.True(t, IsNotFound(apiErr), "404 API错误应该被识别为NotFound")

	// 测试其他错误
	assert.False(t, IsNotFound(errors.New("random error")), "随机错误不应该被识别为NotFound")

	// 测试其他状态码的API错误
	apiErr = &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusBadRequest,
		URL:        "https://example.com/test",
	}
	assert.False(t, IsNotFound(apiErr), "400 API错误不应该被识别为NotFound")
}

// 测试RateLimited错误判断
func TestIsRateLimited(t *testing.T) {
	// 测试直接的ErrRateLimited
	assert.True(t, IsRateLimited(ErrRateLimited), "ErrRateLimited应该被识别为RateLimited")

	// 测试429状态码的API错误
	apiErr := &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusTooManyRequests,
		URL:        "https://example.com/test",
	}
	assert.True(t, IsRateLimited(apiErr), "429 API错误应该被识别为RateLimited")

	// 测试其他错误
	assert.False(t, IsRateLimited(errors.New("random error")), "随机错误不应该被识别为RateLimited")

	// 测试其他状态码的API错误
	apiErr = &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusBadRequest,
		URL:        "https://example.com/test",
	}
	assert.False(t, IsRateLimited(apiErr), "400 API错误不应该被识别为RateLimited")
}

// 测试Unauthorized错误判断
func TestIsUnauthorized(t *testing.T) {
	// 测试直接的ErrUnauthorized
	assert.True(t, IsUnauthorized(ErrUnauthorized), "ErrUnauthorized应该被识别为Unauthorized")

	// 测试401状态码的API错误
	apiErr := &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusUnauthorized,
		URL:        "https://example.com/test",
	}
	assert.True(t, IsUnauthorized(apiErr), "401 API错误应该被识别为Unauthorized")

	// 测试其他错误
	assert.False(t, IsUnauthorized(errors.New("random error")), "随机错误不应该被识别为Unauthorized")

	// 测试其他状态码的API错误
	apiErr = &APIError{
		Cause:      errors.New("test error"),
		StatusCode: http.StatusBadRequest,
		URL:        "https://example.com/test",
	}
	assert.False(t, IsUnauthorized(apiErr), "400 API错误不应该被识别为Unauthorized")
}

// 测试包装错误的类型判断
func TestErrorWrapping(t *testing.T) {
	// 创建一个API错误
	apiErr := &APIError{
		Cause:      ErrNotFound,
		StatusCode: http.StatusNotFound,
		URL:        "https://example.com/test",
	}

	// 再包装一层
	wrappedErr := errors.New("wrapped: " + apiErr.Error())

	// errors.Is应该找不到原始错误，因为没有正确实现
	assert.False(t, errors.Is(wrappedErr, ErrNotFound), "简单包装不应该能识别底层错误")

	// 使用正确的包装方式
	wrappedErr2 := fmt.Errorf("wrapped: %w", apiErr)

	// errors.As应该能够提取API错误
	var extractedAPIErr *APIError
	assert.True(t, errors.As(wrappedErr2, &extractedAPIErr), "errors.As应该能提取API错误")
	assert.Equal(t, http.StatusNotFound, extractedAPIErr.StatusCode, "提取的API错误应该保留状态码")
}

// 测试不同错误类型
func TestErrorTypes(t *testing.T) {
	errorTypes := []error{
		ErrInvalidRequest,
		ErrNotFound,
		ErrServerError,
		ErrRateLimited,
		ErrUnauthorized,
		ErrTimeout,
		ErrNetworkFailure,
	}

	// 确保所有错误类型都是不同的
	for i, err1 := range errorTypes {
		for j, err2 := range errorTypes {
			if i != j {
				assert.NotEqual(t, err1, err2, "错误类型应该是不同的: %v 和 %v", err1, err2)
			}
		}
	}
}
