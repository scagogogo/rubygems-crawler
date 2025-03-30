package repository

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultRetryOptions(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Check default values
	assert.Equal(t, DefaultRetryAttempts, options.MaxAttempts)
	assert.Equal(t, DefaultRetryWaitTime, options.WaitTime)
	assert.Equal(t, DefaultRetryMaxWaitTime, options.MaxWaitTime)
	assert.True(t, options.UseExponentialBackoff)
	assert.NotNil(t, options.ShouldRetry)

	// Test the default shouldRetry function with nil response and error
	assert.True(t, options.ShouldRetry(nil, assert.AnError))

	// Test the default shouldRetry function with various status codes
	// Should retry on 429, 500, 502, 503, 504
	statusCodes := map[int]bool{
		http.StatusOK:                  false, // 200
		http.StatusBadRequest:          false, // 400
		http.StatusUnauthorized:        false, // 401
		http.StatusForbidden:           false, // 403
		http.StatusNotFound:            false, // 404
		http.StatusTooManyRequests:     true,  // 429
		http.StatusInternalServerError: true,  // 500
		http.StatusBadGateway:          true,  // 502
		http.StatusServiceUnavailable:  true,  // 503
		http.StatusGatewayTimeout:      true,  // 504
	}

	for code, shouldRetry := range statusCodes {
		resp := &http.Response{StatusCode: code}
		assert.Equal(t, shouldRetry, options.ShouldRetry(resp, nil), "Status code %d should return %v", code, shouldRetry)
	}
}

func TestRetryOptions_WithMaxAttempts(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Test fluent interface
	result := options.WithMaxAttempts(10)
	assert.Same(t, options, result)

	// Verify value was set
	assert.Equal(t, 10, options.MaxAttempts)
}

func TestRetryOptions_WithWaitTime(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Test fluent interface
	waitTime := 5 * time.Second
	result := options.WithWaitTime(waitTime)
	assert.Same(t, options, result)

	// Verify value was set
	assert.Equal(t, waitTime, options.WaitTime)
}

func TestRetryOptions_WithMaxWaitTime(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Test fluent interface
	maxWaitTime := 60 * time.Second
	result := options.WithMaxWaitTime(maxWaitTime)
	assert.Same(t, options, result)

	// Verify value was set
	assert.Equal(t, maxWaitTime, options.MaxWaitTime)
}

func TestRetryOptions_WithExponentialBackoff(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Test fluent interface with disabling exponential backoff
	result := options.WithExponentialBackoff(false)
	assert.Same(t, options, result)

	// Verify value was set
	assert.False(t, options.UseExponentialBackoff)

	// Test enabling it again
	options.WithExponentialBackoff(true)
	assert.True(t, options.UseExponentialBackoff)
}

func TestRetryOptions_WithShouldRetry(t *testing.T) {
	options := NewDefaultRetryOptions()

	// Create a custom retry function that only retries on 500 errors
	customShouldRetry := func(resp *http.Response, err error) bool {
		if err != nil {
			return true
		}
		if resp != nil && resp.StatusCode == http.StatusInternalServerError {
			return true
		}
		return false
	}

	// Test fluent interface
	result := options.WithShouldRetry(customShouldRetry)
	assert.Same(t, options, result)

	// Verify function was set by testing it
	assert.True(t, options.ShouldRetry(nil, assert.AnError))
	assert.True(t, options.ShouldRetry(&http.Response{StatusCode: http.StatusInternalServerError}, nil))
	assert.False(t, options.ShouldRetry(&http.Response{StatusCode: http.StatusBadGateway}, nil))
	assert.False(t, options.ShouldRetry(&http.Response{StatusCode: http.StatusServiceUnavailable}, nil))
}
