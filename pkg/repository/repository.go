// Package repository 提供了RubyGems API的Go客户端实现
// 它支持官方源和多个国内镜像源，并具有错误处理、重试机制和缓存支持
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/rubygems-crawler/pkg/models"
)

// https://guides.rubygems.org/api/v2/rubygems/[GEM%20NAME]/versions/[VERSION%20NUMBER].(json%7Cyaml)

// Repository 定义了RubyGems API操作的接口
// 它包含了所有与RubyGems交互的核心方法
type Repository interface {
	// GetPackage 通过包名获取包的详细信息
	// 包信息包括名称、版本、作者、下载量、主页URL等
	// 如果包不存在，将返回NotFound错误
	GetPackage(ctx context.Context, gemName string) (*models.PackageInformation, error)

	// Search 根据查询字符串搜索包
	// query参数可以是包名的一部分
	// 返回的结果按照相关性和流行度排序
	// 如果找不到匹配的包，将返回空切片而不是错误
	Search(ctx context.Context, query string, page int) ([]*models.PackageInformation, error)

	// GetGemVersions 获取指定包的所有版本信息
	// 返回的版本按照发布时间降序排列（最新的版本在前）
	// 如果包不存在，将返回空切片而不是错误
	GetGemVersions(ctx context.Context, gemName string) ([]*models.Version, error)

	// GetGemLatestVersion 获取给定包的最新版本
	// GET - /api/v1/versions/[GEM NAME]/latest.json
	GetGemLatestVersion(ctx context.Context, gemName string) (*models.LatestVersion, error)

	// GetTimeFrameVersions 获取特定时间段内的版本信息
	// GET - /api/v1/timeframe_versions.json
	// 时间格式样例: 2019-01-18T21:24:29Z
	GetTimeFrameVersions(ctx context.Context, from, to time.Time) ([]*models.Version, error)

	// Downloads 获取这个仓库中的包总共被下载了多少次
	// GET - /api/v1/downloads.(json|yaml)
	// Returns an object containing the total number of downloads on RubyGems.
	Downloads(ctx context.Context) (*models.RepositoryDownloadCount, error)

	// VersionDownloads 获取给定的包的给定版本总共被下载了多少次
	// GET - /api/v1/downloads/[GEM NAME]-[GEM VERSION].(json|yaml)
	VersionDownloads(ctx context.Context, gemName, gemVersion string) (*models.VersionDownloadCount, error)

	// GetDependencies 获取指定gem包的依赖
	// GET - /api/v1/dependencies?gems=[COMMA DELIMITED GEM NAMES]
	GetDependencies(ctx context.Context, gemsNames ...string) ([]*models.DependencyInfo, error)

	// LatestGems 获取仓库上最新发布的gem包
	// GET - /api/v1/activity/latest.json
	LatestGems(ctx context.Context) ([]*models.PackageInformation, error)

	// GetReverseDependencies 获取依赖于指定gem包的所有包
	// GET - /api/v1/gems/[GEM NAME]/reverse_dependencies.json
	GetReverseDependencies(ctx context.Context, gemName string) ([]string, error)

	// BulkGetPackages 批量获取多个包的信息
	// 并发执行GetPackage请求，提高大规模数据获取效率
	BulkGetPackages(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[*models.PackageInformation]

	// BulkGetVersions 批量获取多个包的版本信息
	// 并发执行GetGemVersions请求，提高大规模数据获取效率
	BulkGetVersions(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]*models.Version]

	// BulkGetDependencies 批量获取多个包的依赖信息
	// 并发执行GetDependencies请求，提高大规模数据获取效率
	BulkGetDependencies(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]*models.DependencyInfo]

	// BulkGetReverseDependencies 批量获取多个包的反向依赖信息
	// 并发执行GetReverseDependencies请求，提高大规模数据获取效率
	BulkGetReverseDependencies(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]string]
}

type RepositoryImpl struct {
	options *Options
}

// NewRepository 创建一个仓库，gem都是存放在仓库中的
func NewRepository(options ...*Options) *RepositoryImpl {
	if len(options) == 0 {
		options = append(options, NewOptions())
	}
	return &RepositoryImpl{
		options: options[0],
	}
}

// GetPackage 获取gem包的基础信息
// GetPackage GET - /api/v1/gems/[GEM NAME].(json|yaml)
func (x *RepositoryImpl) GetPackage(ctx context.Context, gemName string) (*models.PackageInformation, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/gems/%s.json", x.options.ServerURL, gemName)
	return getJson[*models.PackageInformation](ctx, x, targetUrl)
}

// Search 在整个仓库中搜索符合条件的包，使用page参数翻页，如果响应列表为空则说明翻到了尾页
// GET - /api/v1/search.(json|yaml)?query=[YOUR QUERY]
func (x *RepositoryImpl) Search(ctx context.Context, query string, page int) ([]*models.PackageInformation, error) {
	if page <= 0 {
		page = 1
	}
	targetUrl := fmt.Sprintf("%s/api/v1/search.json?query=%s&page=%d", x.options.ServerURL, query, page)
	return getJson[[]*models.PackageInformation](ctx, x, targetUrl)
}

// GetGemVersions 获取指定的gem包的所有版本都有哪些
// GET - /api/v1/versions/[GEM NAME].(json|yaml)
func (x *RepositoryImpl) GetGemVersions(ctx context.Context, gemName string) ([]*models.Version, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/versions/%s.json", x.options.ServerURL, gemName)
	return getJson[[]*models.Version](ctx, x, targetUrl)
}

// GetGemLatestVersion 获取给定包的最新版本
// GET - /api/v1/versions/[GEM NAME]/latest.json
func (x *RepositoryImpl) GetGemLatestVersion(ctx context.Context, gemName string) (*models.LatestVersion, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/versions/%s/latest.json", x.options.ServerURL, gemName)
	return getJson[*models.LatestVersion](ctx, x, targetUrl)
}

// GetTimeFrameVersions 获取特定时间段内的版本信息
// GET - /api/v1/timeframe_versions.json
// 时间格式样例: 2019-01-18T21:24:29Z
func (x *RepositoryImpl) GetTimeFrameVersions(ctx context.Context, from, to time.Time) ([]*models.Version, error) {
	// 格式化时间为RFC3339格式
	fromStr := from.Format(time.RFC3339)
	toStr := to.Format(time.RFC3339)
	targetUrl := fmt.Sprintf("%s/api/v1/timeframe_versions.json?from=%s&to=%s", x.options.ServerURL, fromStr, toStr)
	return getJson[[]*models.Version](ctx, x, targetUrl)
}

// Downloads 获取这个仓库中的包总共被下载了多少次
// GET - /api/v1/downloads.(json|yaml)
// Returns an object containing the total number of downloads on RubyGems.
func (x *RepositoryImpl) Downloads(ctx context.Context) (*models.RepositoryDownloadCount, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/downloads.json", x.options.ServerURL)
	return getJson[*models.RepositoryDownloadCount](ctx, x, targetUrl)
}

// VersionDownloads 获取给定的包的给定版本总共被下载了多少次
// GET - /api/v1/downloads/[GEM NAME]-[GEM VERSION].(json|yaml)
func (x *RepositoryImpl) VersionDownloads(ctx context.Context, gemName, gemVersion string) (*models.VersionDownloadCount, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/downloads/%s-%s.json", x.options.ServerURL, gemName, gemVersion)
	return getJson[*models.VersionDownloadCount](ctx, x, targetUrl)
}

// GetDependencies 获取指定gem包的依赖
// GET - /api/v1/dependencies?gems=[COMMA DELIMITED GEM NAMES]
func (x *RepositoryImpl) GetDependencies(ctx context.Context, gemsNames ...string) ([]*models.DependencyInfo, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/dependencies?gems=%s", x.options.ServerURL, strings.Join(gemsNames, ","))
	return getJson[[]*models.DependencyInfo](ctx, x, targetUrl)
}

// LatestGems 获取仓库上最新发布的gem包
// GET - /api/v1/activity/latest.json
func (x *RepositoryImpl) LatestGems(ctx context.Context) ([]*models.PackageInformation, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/activity/latest.json", x.options.ServerURL)
	return getJson[[]*models.PackageInformation](ctx, x, targetUrl)
}

// GetReverseDependencies 获取依赖于指定gem包的所有包
// GET - /api/v1/gems/[GEM NAME]/reverse_dependencies.json
func (x *RepositoryImpl) GetReverseDependencies(ctx context.Context, gemName string) ([]string, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/gems/%s/reverse_dependencies.json", x.options.ServerURL, gemName)
	return getJson[[]string](ctx, x, targetUrl)
}

func getJson[T any](ctx context.Context, repository *RepositoryImpl, targetUrl string) (T, error) {
	bytes, err := repository.getBytes(ctx, targetUrl)
	if err != nil {
		var zero T
		return zero, err
	}
	return unmarshalJson[T](bytes)
}

func unmarshalJson[T any](bytes []byte) (T, error) {
	var r T
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		var zero T
		return zero, err
	}
	return r, nil
}

// 内部使用统一的方法来请求
func (x *RepositoryImpl) getBytes(ctx context.Context, targetUrl string) ([]byte, error) {
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())

	// 设置代理
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}

	// 设置Token认证
	if x.options.Token != "" {
		// 使用匿名函数方式设置HTTP头
		options.AppendRequestSetting(func(client *http.Client, request *http.Request) error {
			request.Header.Set("Authorization", "Bearer "+x.options.Token)
			return nil
		})
	}

	// 如果启用了重试，使用带重试的请求
	if x.options.RetryOptions != nil {
		return SendRequestWithRetry(ctx, options, x.options.RetryOptions)
	}

	// 否则直接发送请求
	return requests.SendRequest[any, []byte](ctx, options)
}
