package repository

import (
	"context"
	"strings"
	"time"

	"github.com/scagogogo/rubygems-crawler/pkg/cache"
	"github.com/scagogogo/rubygems-crawler/pkg/models"
)

const (
	// DefaultCacheExpiration 默认缓存过期时间 (10分钟)
	DefaultCacheExpiration = 10 * time.Minute

	// DefaultCleanupInterval 默认清理间隔 (1小时)
	DefaultCleanupInterval = 1 * time.Hour
)

// CachedRepository 是带缓存功能的仓库包装器
// 它实现了Repository接口，可以无缝替代基础仓库
// 通过缓存API响应数据，减少重复请求，提高性能
type CachedRepository struct {
	repo          Repository    // 底层仓库实现
	defaultTTL    time.Duration // 默认缓存过期时间
	cache         cache.Cache   // 缓存实现
	stopCleanupCh chan struct{} // 用于停止清理协程的通道
}

// NewCachedRepository 创建一个新的带缓存的仓库实例
// 参数：
//   - repo: 底层仓库实现
//   - ttl: 默认缓存过期时间，所有缓存项的生存时间
//   - cache: 缓存实现，如果为nil，将创建一个新的内存缓存
func NewCachedRepository(repo Repository, ttl time.Duration, cache cache.Cache) *CachedRepository {
	if cache == nil {
		// 如果未提供缓存，创建一个内存缓存，默认清理间隔为缓存时间的两倍
		cache = cache.NewMemoryCache(ttl, ttl*2)
	}

	return &CachedRepository{
		repo:          repo,
		defaultTTL:    ttl,
		cache:         cache,
		stopCleanupCh: make(chan struct{}),
	}
}

// GetPackage 通过缓存获取包信息
// 优先从缓存获取，缓存未命中时调用底层仓库方法并缓存结果
func (c *CachedRepository) GetPackage(ctx context.Context, gemName string) (*models.PackageInformation, error) {
	cacheKey := "package:" + gemName

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if pkg, ok := cachedValue.(*models.PackageInformation); ok {
			return pkg, nil
		}
	}

	// 缓存未命中，调用底层仓库
	pkg, err := c.repo.GetPackage(ctx, gemName)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	c.cache.SetWithExpiration(cacheKey, pkg, c.defaultTTL)
	return pkg, nil
}

// Search 通过缓存执行搜索操作
// 由于搜索结果可能随时间变化，搜索结果的缓存时间较短
func (c *CachedRepository) Search(ctx context.Context, query string, page int) ([]*models.PackageInformation, error) {
	cacheKey := "search:" + query + ":" + string(page)

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if results, ok := cachedValue.([]*models.PackageInformation); ok {
			return results, nil
		}
	}

	// 缓存未命中，调用底层仓库
	results, err := c.repo.Search(ctx, query, page)
	if err != nil {
		return nil, err
	}

	// 搜索结果缓存时间较短，使用默认TTL的一半
	c.cache.SetWithExpiration(cacheKey, results, c.defaultTTL/2)
	return results, nil
}

// GetGemVersions 通过缓存获取包的版本列表
// 版本列表相对稳定，使用默认缓存时间
func (c *CachedRepository) GetGemVersions(ctx context.Context, gemName string) ([]*models.Version, error) {
	cacheKey := "versions:" + gemName

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if versions, ok := cachedValue.([]*models.Version); ok {
			return versions, nil
		}
	}

	// 缓存未命中，调用底层仓库
	versions, err := c.repo.GetGemVersions(ctx, gemName)
	if err != nil {
		return nil, err
	}

	c.cache.SetWithExpiration(cacheKey, versions, c.defaultTTL)
	return versions, nil
}

// GetGemLatestVersion 通过缓存获取包的最新版本
// 由于最新版本可能更新频繁，缓存时间较短
func (c *CachedRepository) GetGemLatestVersion(ctx context.Context, gemName string) (*models.LatestVersion, error) {
	cacheKey := "latest_version:" + gemName

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if version, ok := cachedValue.(*models.LatestVersion); ok {
			return version, nil
		}
	}

	// 缓存未命中，调用底层仓库
	version, err := c.repo.GetGemLatestVersion(ctx, gemName)
	if err != nil {
		return nil, err
	}

	// 最新版本缓存时间较短
	c.cache.SetWithExpiration(cacheKey, version, c.defaultTTL/2)
	return version, nil
}

// GetTimeFrameVersions 通过缓存获取时间段内的版本
// 时间段查询结果相对稳定，使用默认缓存时间
func (c *CachedRepository) GetTimeFrameVersions(ctx context.Context, from, to time.Time) ([]*models.Version, error) {
	cacheKey := "timeframe:" + from.Format(time.RFC3339) + ":" + to.Format(time.RFC3339)

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if versions, ok := cachedValue.([]*models.Version); ok {
			return versions, nil
		}
	}

	// 缓存未命中，调用底层仓库
	versions, err := c.repo.GetTimeFrameVersions(ctx, from, to)
	if err != nil {
		return nil, err
	}

	c.cache.SetWithExpiration(cacheKey, versions, c.defaultTTL)
	return versions, nil
}

// Downloads 通过缓存获取仓库下载统计
// 下载统计变化较频繁，使用较短的缓存时间
func (c *CachedRepository) Downloads(ctx context.Context) (*models.RepositoryDownloadCount, error) {
	cacheKey := "downloads"

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if downloads, ok := cachedValue.(*models.RepositoryDownloadCount); ok {
			return downloads, nil
		}
	}

	// 缓存未命中，调用底层仓库
	downloads, err := c.repo.Downloads(ctx)
	if err != nil {
		return nil, err
	}

	// 下载统计缓存时间较短
	c.cache.SetWithExpiration(cacheKey, downloads, c.defaultTTL/2)
	return downloads, nil
}

// VersionDownloads 通过缓存获取特定版本的下载统计
// 版本下载统计变化较频繁，使用较短的缓存时间
func (c *CachedRepository) VersionDownloads(ctx context.Context, gemName, gemVersion string) (*models.VersionDownloadCount, error) {
	cacheKey := "version_downloads:" + gemName + ":" + gemVersion

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if downloads, ok := cachedValue.(*models.VersionDownloadCount); ok {
			return downloads, nil
		}
	}

	// 缓存未命中，调用底层仓库
	downloads, err := c.repo.VersionDownloads(ctx, gemName, gemVersion)
	if err != nil {
		return nil, err
	}

	// 版本下载统计缓存时间较短
	c.cache.SetWithExpiration(cacheKey, downloads, c.defaultTTL/2)
	return downloads, nil
}

// GetDependencies 通过缓存获取包的依赖关系
// 依赖关系相对稳定，使用默认缓存时间
func (c *CachedRepository) GetDependencies(ctx context.Context, gemNames ...string) ([]*models.DependencyInfo, error) {
	// 对于多个包名，使用连接字符串作为缓存键
	cacheKey := "dependencies:" + strings.Join(gemNames, ",")

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if deps, ok := cachedValue.([]*models.DependencyInfo); ok {
			return deps, nil
		}
	}

	// 缓存未命中，调用底层仓库
	deps, err := c.repo.GetDependencies(ctx, gemNames...)
	if err != nil {
		return nil, err
	}

	c.cache.SetWithExpiration(cacheKey, deps, c.defaultTTL)
	return deps, nil
}

// LatestGems 通过缓存获取最新的gem包列表
// 最新列表变化频繁，使用较短的缓存时间
func (c *CachedRepository) LatestGems(ctx context.Context) ([]*models.PackageInformation, error) {
	cacheKey := "latest_gems"

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if gems, ok := cachedValue.([]*models.PackageInformation); ok {
			return gems, nil
		}
	}

	// 缓存未命中，调用底层仓库
	gems, err := c.repo.LatestGems(ctx)
	if err != nil {
		return nil, err
	}

	// 最新列表缓存时间较短
	c.cache.SetWithExpiration(cacheKey, gems, c.defaultTTL/4)
	return gems, nil
}

// GetReverseDependencies 通过缓存获取包的反向依赖
// 反向依赖相对稳定，使用默认缓存时间
func (c *CachedRepository) GetReverseDependencies(ctx context.Context, gemName string) ([]string, error) {
	cacheKey := "reverse_dependencies:" + gemName

	// 尝试从缓存获取
	if cachedValue, ok := c.cache.Get(cacheKey); ok {
		if deps, ok := cachedValue.([]string); ok {
			return deps, nil
		}
	}

	// 缓存未命中，调用底层仓库
	deps, err := c.repo.GetReverseDependencies(ctx, gemName)
	if err != nil {
		return nil, err
	}

	c.cache.SetWithExpiration(cacheKey, deps, c.defaultTTL)
	return deps, nil
}

// Close 关闭缓存仓库，释放资源
// 在仓库不再使用时应调用此方法
func (c *CachedRepository) Close() {
	close(c.stopCleanupCh)
}

// ClearCache 清空缓存
// 可在需要强制刷新数据时调用
func (c *CachedRepository) ClearCache() {
	c.cache.Clear()
}

// GetCacheStats 获取当前缓存统计信息
// 返回当前缓存中的项目数量
func (c *CachedRepository) GetCacheStats() int {
	return c.cache.ItemCount()
}
