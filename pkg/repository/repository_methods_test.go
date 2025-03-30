package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试获取时间段内的版本
func TestRepository_GetTimeFrameVersions(t *testing.T) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建仓库实例
	repo := NewRepository()

	// 设置时间范围，选择最近24小时
	to := time.Now()
	from := to.Add(-24 * time.Hour)

	// 获取时间段内的版本
	versions, err := repo.GetTimeFrameVersions(ctx, from, to)

	// 验证结果
	assert.NoError(t, err, "获取时间段内的版本不应该返回错误")
	assert.NotNil(t, versions, "返回的版本列表不应为nil")

	// 如果没有返回版本，只是说明这段时间内没有版本发布，不算错误
	if len(versions) > 0 {
		// 验证版本的创建时间是否在指定范围内
		for _, version := range versions {
			assert.True(t, version.CreatedAt.After(from) && version.CreatedAt.Before(to.Add(time.Minute)),
				"版本创建时间应该在指定范围内: %v-%v, 实际: %v", from, to, version.CreatedAt)
			assert.NotEmpty(t, version.Number, "版本号不能为空")
		}
	}
}

// 测试获取包的依赖
func TestRepository_GetDependencies(t *testing.T) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建仓库实例
	repo := NewRepository()

	// 测试单个包依赖
	t.Run("单个包依赖", func(t *testing.T) {
		// 选择一个常见包，它应该有依赖
		dependencies, err := repo.GetDependencies(ctx, "rails")

		assert.NoError(t, err, "获取依赖不应返回错误")
		assert.NotNil(t, dependencies, "依赖列表不应为nil")
		assert.NotEmpty(t, dependencies, "Rails应该有依赖")

		// 检查依赖项的字段
		for _, dep := range dependencies {
			assert.NotEmpty(t, dep.Name, "依赖项名不能为空")
			assert.NotEmpty(t, dep.Requirements, "依赖要求不能为空")
		}
	})

	// 测试多个包依赖
	t.Run("多个包依赖", func(t *testing.T) {
		// 选择几个常见包
		dependencies, err := repo.GetDependencies(ctx, "rails", "rack", "nokogiri")

		assert.NoError(t, err, "获取多个包依赖不应返回错误")
		assert.NotNil(t, dependencies, "依赖列表不应为nil")
		assert.NotEmpty(t, dependencies, "这些包应该有依赖")
	})

	// 测试获取不存在包的依赖
	t.Run("不存在的包", func(t *testing.T) {
		// 使用一个极大概率不存在的包名
		dependencies, err := repo.GetDependencies(ctx, "non_existent_package_xyz_123")

		// 这里应该返回空列表而不是错误
		assert.NoError(t, err, "获取不存在包的依赖应返回空列表，不是错误")
		assert.Empty(t, dependencies, "不存在的包不应有依赖")
	})
}

// 测试获取包的反向依赖
func TestRepository_GetReverseDependencies(t *testing.T) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建仓库实例
	repo := NewRepository()

	// 测试常用包的反向依赖
	t.Run("常用包反向依赖", func(t *testing.T) {
		// 选择一个基础包，它应该被很多其他包依赖
		dependencies, err := repo.GetReverseDependencies(ctx, "rack")

		assert.NoError(t, err, "获取反向依赖不应返回错误")
		assert.NotNil(t, dependencies, "反向依赖列表不应为nil")
		assert.NotEmpty(t, dependencies, "Rack应该有反向依赖")
	})

	// 测试获取不存在包的反向依赖
	t.Run("不存在的包", func(t *testing.T) {
		// 使用一个极大概率不存在的包名
		dependencies, err := repo.GetReverseDependencies(ctx, "non_existent_package_xyz_123")

		// 这里应该返回空列表而不是错误
		assert.NoError(t, err, "获取不存在包的反向依赖应返回空列表，不是错误")
		assert.Empty(t, dependencies, "不存在的包不应有反向依赖")
	})

	// 测试新包的反向依赖
	t.Run("新包反向依赖", func(t *testing.T) {
		// 先获取最新的包列表
		latestGems, err := repo.LatestGems(ctx)
		if err != nil || len(latestGems) == 0 {
			t.Skip("无法获取最新包列表")
			return
		}

		// 选择最新发布的包，它可能没有反向依赖
		dependencies, err := repo.GetReverseDependencies(ctx, latestGems[0].Name)

		// 不关心是否有依赖，但不应该出错
		assert.NoError(t, err, "获取新包的反向依赖不应返回错误")
		assert.NotNil(t, dependencies, "反向依赖列表不应为nil")
	})
}

// 测试使用不同的镜像源
func TestDifferentMirrors(t *testing.T) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建不同镜像源的仓库
	repos := map[string]*Repository{
		"默认":        NewRepository(),
		"RubyChina": NewRubyChinaRepository(),
		"TsingHua":  NewTSingHuaRepository(),
		"AliYun":    NewAliYunRepository(),
	}

	// 对每个镜像源进行测试
	for name, repo := range repos {
		t.Run(name, func(t *testing.T) {
			// 测试获取包信息
			pkg, err := repo.GetPackage(ctx, "rails")
			assert.NoError(t, err, "%s: 获取包信息失败", name)
			assert.NotNil(t, pkg, "%s: 包信息为nil", name)
			assert.Equal(t, "rails", pkg.Name, "%s: 包名不匹配", name)
		})
	}
}

// 测试代理设置 (这个测试需要有效的代理，不在CI环境中运行)
func TestProxySetting(t *testing.T) {
	// 跳过CI环境中的测试
	if testing.Short() {
		t.Skip("在短模式下跳过代理测试")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建带代理的仓库选项
	// 注意: 这是一个示例，需要替换为实际可用的代理
	proxyURL := ""
	if proxyURL == "" {
		t.Skip("未设置代理URL，跳过测试")
	}

	options := NewOptions().SetProxy(proxyURL)
	repo := NewRepository(options)

	// 测试基本功能是否正常
	pkg, err := repo.GetPackage(ctx, "rails")
	assert.NoError(t, err, "通过代理获取包信息失败")
	assert.NotNil(t, pkg, "包信息为nil")
	assert.Equal(t, "rails", pkg.Name, "包名不匹配")
}

// 测试Token设置 (这个测试需要有效的Token，不在CI环境中运行)
func TestTokenSetting(t *testing.T) {
	// 跳过CI环境中的测试
	if testing.Short() {
		t.Skip("在短模式下跳过Token测试")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建带Token的仓库选项
	// 注意: 这是一个示例，需要替换为实际可用的Token
	token := ""
	if token == "" {
		t.Skip("未设置Token，跳过测试")
	}

	options := NewOptions().SetToken(token)
	repo := NewRepository(options)

	// 测试基本功能是否正常
	pkg, err := repo.GetPackage(ctx, "rails")
	assert.NoError(t, err, "使用Token获取包信息失败")
	assert.NotNil(t, pkg, "包信息为nil")
	assert.Equal(t, "rails", pkg.Name, "包名不匹配")
}
