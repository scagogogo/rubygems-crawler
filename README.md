# RubyGems爬虫 (RubyGems Crawler)

这是一个Go语言编写的RubyGems API客户端，用于获取RubyGems.org的包信息。它支持多种API操作、国内镜像源、错误处理和缓存机制。

[![GoDoc](https://godoc.org/github.com/scagogogo/rubygems-crawler?status.svg)](https://godoc.org/github.com/scagogogo/rubygems-crawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/rubygems-crawler)](https://goreportcard.com/report/github.com/scagogogo/rubygems-crawler)

## 功能特点

- 支持RubyGems API v1/v2的全部主要功能
- 提供多个国内镜像源支持（Ruby China、清华大学、阿里云）
- 智能错误处理和自动重试机制
- HTTP代理支持和API Token认证
- 内存缓存机制，支持自定义过期时间
- 批量并发请求功能，提高大规模数据获取效率
- 完整的单元测试，保证代码质量
- 完整的命令行工具，支持JSON格式输出

## 速率限制

RubyGems.org API有速率限制，详情参考官方文档:

```
https://guides.rubygems.org/rubygems-org-rate-limits/
```

使用Token认证可以提高API请求配额。

## 安装

```bash
go get github.com/scagogogo/rubygems-crawler
```

## 快速开始

### 基本使用示例

```go
package main

import (
	"context"
	"fmt"
	"github.com/scagogogo/rubygems-crawler/pkg/repository"
)

func main() {
	// 创建默认仓库客户端
	repo := repository.NewRepository()
	
	// 获取特定gem包信息
	pkg, err := repo.GetPackage(context.Background(), "rails")
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Rails版本: %s\n", pkg.Version)
	fmt.Printf("下载量: %d\n", pkg.Downloads)
	fmt.Printf("作者: %s\n", pkg.Authors)
}
```

### 使用国内镜像源

```go
// 使用Ruby中国镜像源
repo := repository.NewRubyChinaRepository()

// 或使用清华大学镜像源
// repo := repository.NewTSingHuaRepository()

// 或使用阿里云镜像源
// repo := repository.NewAliYunRepository()
```

### 使用缓存机制

```go
// 创建内存缓存，设置默认过期时间和清理间隔
memCache := cache.NewMemoryCache(10*time.Minute, 30*time.Minute)

// 创建缓存仓库包装器，默认缓存时间5分钟
cachedRepo := repository.NewCachedRepository(repo, 5*time.Minute, memCache)

// 使用缓存仓库像普通仓库一样进行操作
pkg, err := cachedRepo.GetPackage(context.Background(), "rails")

// 清空缓存
cachedRepo.ClearCache()

// 关闭缓存
defer cachedRepo.Close()
```

### 批量并发请求

```go
// 定义要批量获取的gem包列表
gems := []string{"rails", "rack", "activesupport", "rake", "bundler"}

// 设置批量操作选项，最大并发数为5
options := repository.NewBulkOptions().WithMaxConcurrency(5)

// 批量获取包信息
results := repo.BulkGetPackages(ctx, gems, options)

// 处理结果
for _, result := range results {
    if result.Error != nil {
        fmt.Printf("获取 %s 失败: %v\n", result.Key, result.Error)
        continue
    }
    
    pkg := result.Value
    fmt.Printf("包名: %s, 版本: %s\n", pkg.Name, pkg.Version)
}
```

### 使用Token认证

```go
// 设置Token
options := repository.NewOptions().SetToken("your-api-token")
repo := repository.NewRepository(options)
```

### 使用代理

```go
// 设置HTTP代理
options := repository.NewOptions().SetProxy("http://127.0.0.1:7890")
repo := repository.NewRepository(options)
```

### 自定义重试策略

```go
// 配置重试策略
retryOptions := repository.NewDefaultRetryOptions().
	WithMaxAttempts(5).
	WithWaitTime(2 * time.Second).
	WithExponentialBackoff(true)
	
options := repository.NewOptions().SetRetryOptions(retryOptions)
repo := repository.NewRepository(options)
```

### 错误处理

```go
pkg, err := repo.GetPackage(ctx, "non-existent-package")
if err != nil {
    if repository.IsNotFound(err) {
        fmt.Println("包不存在")
    } else if repository.IsRateLimited(err) {
        fmt.Println("API请求被限流")
    } else if repository.IsUnauthorized(err) {
        fmt.Println("认证失败")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 命令行工具

项目提供了命令行工具，可以直接在终端使用：

```bash
# 获取包信息
rubygems-cli -get -gem rails

# 搜索包
rubygems-cli -search -query rails -limit 10

# 获取版本列表
rubygems-cli -versions -gem rails -limit 20

# 获取依赖信息
rubygems-cli -deps -gem rails

# 获取反向依赖
rubygems-cli -rdeps -gem rails -limit 50

# 使用JSON格式输出
rubygems-cli -get -gem rails -json

# 使用镜像源
rubygems-cli -get -gem rails -mirror ruby-china

# 启用缓存
rubygems-cli -get -gem rails -cache
```

## 项目结构

```
├── cmd/                  # 命令行工具
│   └── rubygems/         # RubyGems命令行客户端
├── examples/             # 使用示例
│   ├── basic_usage.go    # 基本使用示例
│   ├── bulk/             # 批量操作示例
│   └── cache/            # 缓存使用示例
├── pkg/                  # 项目核心包
│   ├── cache/            # 缓存实现
│   ├── models/           # 数据模型
│   └── repository/       # 仓库实现
└── tests/                # 测试目录
    └── integration/      # 集成测试
```

## 单元测试

项目包含了全面的单元测试，覆盖了所有核心功能：

- **模型测试**: 测试了所有数据模型的JSON序列化和反序列化
- **缓存测试**: 测试了缓存的存储、过期、清理和关闭等功能
- **仓库测试**: 测试了仓库选项、重试机制和镜像源配置等

运行测试：

```bash
# 运行所有测试
go test -v ./...

# 运行特定包的测试
go test -v ./pkg/models/...
go test -v ./pkg/cache/...
go test -v ./pkg/repository/...

# 运行带网络访问的测试（默认跳过）
go test -v -run TestLiveAPI ./pkg/repository/...
```

## API参考

详细的API文档请参考代码注释和[RubyGems API文档](https://guides.rubygems.org/rubygems-org-api-v2/)。

### 主要功能列表

#### Repository接口

- `GetPackage(ctx, gemName)`: 获取包详细信息
- `Search(ctx, query, page)`: 搜索包
- `GetGemVersions(ctx, gemName)`: 获取包的所有版本
- `GetGemLatestVersion(ctx, gemName)`: 获取包的最新版本
- `GetTimeFrameVersions(ctx, from, to)`: 获取特定时间段内的版本
- `Downloads(ctx)`: 获取总下载统计
- `VersionDownloads(ctx, gemName, gemVersion)`: 获取特定版本的下载统计
- `GetDependencies(ctx, gemsNames...)`: 获取包的依赖
- `LatestGems(ctx)`: 获取最新发布的包
- `GetReverseDependencies(ctx, gemName)`: 获取依赖于特定包的所有包

#### Cache接口

- `Get(key)`: 获取缓存值
- `Set(key, value)`: 设置缓存值
- `SetWithExpiration(key, value, duration)`: 设置带过期时间的缓存值
- `Delete(key)`: 删除缓存项
- `Clear()`: 清空缓存
- `Count()`: 获取缓存项数量
- `Close()`: 关闭缓存

## 贡献

欢迎提交PR和Issue！在提交PR前，请确保您的代码：

1. 通过了所有测试（`go test ./...`）
2. 添加了新功能的单元测试（如适用）
3. 更新了文档（如适用）
4. 符合项目的代码风格（使用`gofmt`格式化）

## 开发计划

- [ ] 添加更多镜像源支持
- [ ] 提高测试覆盖率到90%+
- [ ] 添加性能测试和基准测试
- [ ] 添加内容下载功能
- [ ] 实现分布式爬虫系统

## 许可证

MIT

## 参考资料

- [RubyGems API v2文档](https://guides.rubygems.org/rubygems-org-api-v2/)
- [RubyGems API v1文档](https://guides.rubygems.org/rubygems-org-api/)
- [RubyGems API速率限制](https://guides.rubygems.org/rubygems-org-rate-limits/)
