package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/scagogogo/rubygems-crawler/pkg/models"
)

// 创建一个模拟的仓库实现用于测试
type mockRepository struct {
	*Repository
	mockPackages map[string]*models.PackageInformation
	mockVersions map[string][]*models.Version
	// 人为延迟，模拟网络请求延迟
	delay time.Duration
	// 人为错误，模拟请求失败
	failOn map[string]error
}

// 创建一个新的模拟仓库
func newMockRepository() *mockRepository {
	repo := &mockRepository{
		Repository:   NewRepository(),
		mockPackages: make(map[string]*models.PackageInformation),
		mockVersions: make(map[string][]*models.Version),
		delay:        10 * time.Millisecond, // 默认10ms延迟
		failOn:       make(map[string]error),
	}

	// 添加一些测试数据
	repo.mockPackages["rails"] = &models.PackageInformation{
		Name:        "rails",
		Version:     "7.0.5",
		Downloads:   1000000,
		HomepageURI: "https://rubyonrails.org",
		Info:        "Ruby on Rails",
	}

	repo.mockPackages["rack"] = &models.PackageInformation{
		Name:        "rack",
		Version:     "2.2.7",
		Downloads:   2000000,
		HomepageURI: "https://github.com/rack/rack",
		Info:        "Rack provides a minimal interface between webservers and Ruby frameworks",
	}

	// 添加一些版本信息
	repo.mockVersions["rails"] = []*models.Version{
		{Number: "7.0.5", CreatedAt: time.Now().Add(-24 * time.Hour)},
		{Number: "7.0.4", CreatedAt: time.Now().Add(-48 * time.Hour)},
	}

	repo.mockVersions["rack"] = []*models.Version{
		{Number: "2.2.7", CreatedAt: time.Now().Add(-24 * time.Hour)},
		{Number: "2.2.6", CreatedAt: time.Now().Add(-48 * time.Hour)},
	}

	return repo
}

// 设置延迟时间
func (m *mockRepository) setDelay(delay time.Duration) *mockRepository {
	m.delay = delay
	return m
}

// 设置特定gem会触发的错误
func (m *mockRepository) setFailOn(gemName string, err error) *mockRepository {
	m.failOn[gemName] = err
	return m
}

// 实现GetPackage方法
func (m *mockRepository) GetPackage(ctx context.Context, gemName string) (*models.PackageInformation, error) {
	// 检查是否应该失败
	if err, ok := m.failOn[gemName]; ok {
		return nil, err
	}

	// 模拟网络延迟
	time.Sleep(m.delay)

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 返回结果
	pkg, ok := m.mockPackages[gemName]
	if !ok {
		return nil, errors.New("gem not found")
	}
	return pkg, nil
}

// 实现GetGemVersions方法
func (m *mockRepository) GetGemVersions(ctx context.Context, gemName string) ([]*models.Version, error) {
	// 检查是否应该失败
	if err, ok := m.failOn[gemName]; ok {
		return nil, err
	}

	// 模拟网络延迟
	time.Sleep(m.delay)

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 返回结果
	versions, ok := m.mockVersions[gemName]
	if !ok {
		return nil, errors.New("gem not found")
	}
	return versions, nil
}

// 测试批量获取包信息
func TestBulkGetPackages(t *testing.T) {
	// 创建模拟仓库
	mockRepo := newMockRepository()

	// 设置一个错误
	mockRepo.setFailOn("not-exist", errors.New("gem not found"))

	// 测试用例
	testCases := []struct {
		name        string
		gemNames    []string
		concurrency int
		timeout     time.Duration
		expectErr   bool
		expectCount int
	}{
		{
			name:        "获取有效包信息",
			gemNames:    []string{"rails", "rack"},
			concurrency: 2,
			timeout:     100 * time.Millisecond,
			expectErr:   false,
			expectCount: 2,
		},
		{
			name:        "包含一个不存在的包",
			gemNames:    []string{"rails", "rack", "not-exist"},
			concurrency: 2,
			timeout:     100 * time.Millisecond,
			expectErr:   true,
			expectCount: 3,
		},
		{
			name:        "超时测试",
			gemNames:    []string{"rails", "rack"},
			concurrency: 1,
			timeout:     5 * time.Millisecond, // 设置很短的超时时间
			expectErr:   true,
			expectCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置上下文和超时时间
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			// 设置并发数
			options := NewBulkOptions().WithMaxConcurrency(tc.concurrency)

			// 执行批量获取
			results := mockRepo.BulkGetPackages(ctx, tc.gemNames, options)

			// 验证结果数量
			if len(results) != tc.expectCount {
				t.Errorf("结果数量不符合预期，期望: %d, 实际: %d", tc.expectCount, len(results))
			}

			// 验证是否有错误
			hasError := false
			for _, result := range results {
				if result.Error != nil {
					hasError = true
					break
				}
			}

			if hasError != tc.expectErr {
				t.Errorf("错误状态不符合预期，期望有错误: %v, 实际: %v", tc.expectErr, hasError)
			}
		})
	}
}

// 测试批量获取版本信息
func TestBulkGetVersions(t *testing.T) {
	// 创建模拟仓库
	mockRepo := newMockRepository()

	// 设置一个错误
	mockRepo.setFailOn("not-exist", errors.New("gem not found"))

	// 测试用例
	testCases := []struct {
		name        string
		gemNames    []string
		concurrency int
		timeout     time.Duration
		expectErr   bool
		expectCount int
	}{
		{
			name:        "获取有效版本信息",
			gemNames:    []string{"rails", "rack"},
			concurrency: 2,
			timeout:     100 * time.Millisecond,
			expectErr:   false,
			expectCount: 2,
		},
		{
			name:        "包含一个不存在的包",
			gemNames:    []string{"rails", "rack", "not-exist"},
			concurrency: 2,
			timeout:     100 * time.Millisecond,
			expectErr:   true,
			expectCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置上下文和超时时间
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			// 设置并发数
			options := NewBulkOptions().WithMaxConcurrency(tc.concurrency)

			// 执行批量获取
			results := mockRepo.BulkGetVersions(ctx, tc.gemNames, options)

			// 验证结果数量
			if len(results) != tc.expectCount {
				t.Errorf("结果数量不符合预期，期望: %d, 实际: %d", tc.expectCount, len(results))
			}

			// 验证是否有错误
			hasError := false
			for _, result := range results {
				if result.Error != nil {
					hasError = true
					break
				}
			}

			if hasError != tc.expectErr {
				t.Errorf("错误状态不符合预期，期望有错误: %v, 实际: %v", tc.expectErr, hasError)
			}
		})
	}
}

// 测试选项设置
func TestBulkOptions(t *testing.T) {
	// 测试默认选项
	options := NewBulkOptions()
	if options.MaxConcurrency != 10 {
		t.Errorf("默认MaxConcurrency应该是10，实际是: %d", options.MaxConcurrency)
	}
	if options.IgnoreErrors != false {
		t.Errorf("默认IgnoreErrors应该是false，实际是: %v", options.IgnoreErrors)
	}

	// 测试设置MaxConcurrency
	options = options.WithMaxConcurrency(5)
	if options.MaxConcurrency != 5 {
		t.Errorf("设置后MaxConcurrency应该是5，实际是: %d", options.MaxConcurrency)
	}

	// 测试设置无效的MaxConcurrency
	options = options.WithMaxConcurrency(0)
	if options.MaxConcurrency != 5 {
		t.Errorf("设置无效值后MaxConcurrency不应变化，实际是: %d", options.MaxConcurrency)
	}

	// 测试设置IgnoreErrors
	options = options.WithIgnoreErrors(true)
	if options.IgnoreErrors != true {
		t.Errorf("设置后IgnoreErrors应该是true，实际是: %v", options.IgnoreErrors)
	}
}
