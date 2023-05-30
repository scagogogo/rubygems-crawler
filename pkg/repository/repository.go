package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/rubygems-crawler/pkg/model"
)

// https://guides.rubygems.org/api/v2/rubygems/[GEM%20NAME]/versions/[VERSION%20NUMBER].(json%7Cyaml)

type Repository struct {
	options *Options
}

// NewRepository 创建一个仓库，gem都是存放在仓库中的
func NewRepository(options ...*Options) *Repository {
	if len(options) == 0 {
		options = append(options, NewOptions())
	}
	return &Repository{
		options: options[0],
	}
}

// GetPackage 获取gem包的基础信息
// GetPackage GET - /api/v1/gems/[GEM NAME].(json|yaml)
func (x *Repository) GetPackage(ctx context.Context, gemName string) (*model.PackageInformation, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/gems/%s.json", x.options.ServerURL, gemName)
	return getJson[*model.PackageInformation](ctx, x, targetUrl)
}

// Search 在整个仓库中搜索符合条件的包，使用page参数翻页，如果响应列表为空则说明翻到了尾页
// GET - /api/v1/search.(json|yaml)?query=[YOUR QUERY]
func (x *Repository) Search(ctx context.Context, query string, page int) ([]*model.PackageInformation, error) {
	if page <= 0 {
		page = 1
	}
	targetUrl := fmt.Sprintf("%s/api/v1/search.json?query=%s&page=%d", x.options.ServerURL, query, page)
	return getJson[[]*model.PackageInformation](ctx, x, targetUrl)
}

// GetGemVersions 获取指定的gem包的所有版本都有哪些
// GET - /api/v1/versions/[GEM NAME].(json|yaml)
func (x *Repository) GetGemVersions(ctx context.Context, gemName string) ([]*model.Version, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/versions/%s.json", x.options.ServerURL, gemName)
	return getJson[[]*model.Version](ctx, x, targetUrl)
}

// GetGemLatestVersion 获取给定包的最新版本
// GET - /api/v1/versions/[GEM NAME]/latest.json
func (x *Repository) GetGemLatestVersion(ctx context.Context, gemName string) (*model.LatestVersion, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/versions/%s/latest.json", x.options.ServerURL, gemName)
	return getJson[*model.LatestVersion](ctx, x, targetUrl)
}

// TODO
//// GET - /api/v1/timeframe_versions.json
//func (x *Repository) GetTimeFrameVersions(ctx context.Context, from, to time.Time) (*model.LatestVersion, error) {
//	// 2019-01-18T21:24:29Z
//	targetUrl := fmt.Sprintf("%s/api/v1/timeframe_versions.json?from=%s&to=%s", x.options.ServerURL, gemName)
//	return getJson[*model.LatestVersion](ctx, x, targetUrl)
//}

// Downloads 获取这个仓库中的包总共被下载了多少次
// GET - /api/v1/downloads.(json|yaml)
// Returns an object containing the total number of downloads on RubyGems.
func (x *Repository) Downloads(ctx context.Context) (*model.RepositoryDownloadCount, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/downloads.json", x.options.ServerURL)
	return getJson[*model.RepositoryDownloadCount](ctx, x, targetUrl)
}

// VersionDownloads 获取给定的包的给定版本总共被下载了多少次
// GET - /api/v1/downloads/[GEM NAME]-[GEM VERSION].(json|yaml)
func (x *Repository) VersionDownloads(ctx context.Context, gemName, gemVersion string) (*model.VersionDownloadCount, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/downloads/%s-%s.json", x.options.ServerURL, gemName, gemVersion)
	return getJson[*model.VersionDownloadCount](ctx, x, targetUrl)
}

//// GET - /api/v1/dependencies?gems=[COMMA DELIMITED GEM NAMES]
//func (x *Repository) GetDependencies(gemsNames ...string) {
//	targetUrl := fmt.Sprintf("%s/api/v1/dependencies?gems=%s", x.options.ServerURL, strings.Join(gemsNames, ","))
//	bytes, err := x.getBytes(ctx, targetUrl)
//	if err != nil {
//		return nil, err
//	}
//
//}

// LatestGems 获取仓库上最新发布的gem包
// GET - /api/v1/activity/latest.json
func (x *Repository) LatestGems(ctx context.Context) ([]*model.PackageInformation, error) {
	targetUrl := fmt.Sprintf("%s/api/v1/activity/latest.json", x.options.ServerURL)
	return getJson[[]*model.PackageInformation](ctx, x, targetUrl)
}

//// https://rubygems.org/api/v1/gems/[GEM NAME]/reverse_dependencies.json
//func (x *Repository) GetName() {
//
//}

func getJson[T any](ctx context.Context, repository *Repository, targetUrl string) (T, error) {
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
func (x *Repository) getBytes(ctx context.Context, targetUrl string) ([]byte, error) {
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	return requests.SendRequest[any, []byte](ctx, options)
}
