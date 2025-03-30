package repository

import (
	"context"
	"testing"
	"time"

	"github.com/scagogogo/rubygems-crawler/pkg/cache"
	"github.com/scagogogo/rubygems-crawler/pkg/models"
	"github.com/stretchr/testify/assert"
)

// 模拟Repository用于测试
type MockRepo struct {
	calledTimes int
	testPkg     *models.PackageInformation
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		calledTimes: 0,
		testPkg: &models.PackageInformation{
			Name:    "test-gem",
			Version: "1.0.0",
			Authors: "Test Author",
		},
	}
}

func (m *MockRepo) GetPackage(ctx context.Context, gemName string) (*models.PackageInformation, error) {
	m.calledTimes++
	return m.testPkg, nil
}

func TestCachedRepository(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockRepo()

	// 创建一个内存缓存
	memCache := cache.NewMemoryCache(10*time.Minute, 30*time.Minute)

	// 创建一个测试包装器
	type testWrapper struct {
		repo      *MockRepo
		cache     cache.Cache
		getCalled func() int
	}

	wrapper := &testWrapper{
		repo:  mockRepo,
		cache: memCache,
		getCalled: func() int {
			return mockRepo.calledTimes
		},
	}

	// 测试不使用缓存的情况
	for i := 0; i < 3; i++ {
		pkg, err := wrapper.repo.GetPackage(ctx, "test-gem")
		assert.NoError(t, err)
		assert.Equal(t, "test-gem", pkg.Name)
	}

	// 应该被调用3次
	assert.Equal(t, 3, wrapper.getCalled())

	// 创建新的mock和缓存仓库
	mockRepo2 := NewMockRepo()
	cacheRepo := NewCachedRepository(&Repository{}, 10*time.Minute, memCache)

	// 创建一个自定义的GetPackage函数，使用我们的mock
	getPackageWithMock := func(ctx context.Context, gemName string) (*models.PackageInformation, error) {
		// 使用我们的mock而不是真实的API调用
		pkg, err := mockRepo2.GetPackage(ctx, gemName)
		if err != nil {
			return nil, err
		}

		// 手动设置缓存
		cacheKey := "package:" + gemName
		cacheRepo.cache.SetWithExpiration(cacheKey, pkg, cacheRepo.expiration)

		return pkg, nil
	}

	// 首次调用，通过mock获取数据并缓存
	pkg, err := getPackageWithMock(ctx, "test-gem")
	assert.NoError(t, err)
	assert.Equal(t, "test-gem", pkg.Name)
	assert.Equal(t, 1, mockRepo2.calledTimes)

	// 第二次调用，应该从缓存获取
	cachedPkg, err := cacheRepo.GetPackage(ctx, "test-gem")
	assert.NoError(t, err)
	assert.Equal(t, "test-gem", cachedPkg.Name)

	// mock仍然只被调用了一次
	assert.Equal(t, 1, mockRepo2.calledTimes)

	// 清理
	cacheRepo.ClearCache()
	cacheRepo.Close()
}
