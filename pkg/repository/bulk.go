package repository

import (
	"context"
	"sync"

	"github.com/scagogogo/rubygems-crawler/pkg/models"
)

// BulkResult 批量操作的结果
type BulkResult[T any] struct {
	// 操作的键（如gem名称）
	Key string

	// 操作结果
	Value T

	// 操作过程中的错误
	Error error
}

// BulkOptions 批量操作的选项
type BulkOptions struct {
	// 最大并发数，默认为10
	MaxConcurrency int

	// 是否忽略错误继续执行
	IgnoreErrors bool
}

// NewBulkOptions 创建批量操作选项
func NewBulkOptions() *BulkOptions {
	return &BulkOptions{
		MaxConcurrency: 10,
		IgnoreErrors:   false,
	}
}

// WithMaxConcurrency 设置最大并发数
func (o *BulkOptions) WithMaxConcurrency(max int) *BulkOptions {
	if max > 0 {
		o.MaxConcurrency = max
	}
	return o
}

// WithIgnoreErrors 设置是否忽略错误
func (o *BulkOptions) WithIgnoreErrors(ignore bool) *BulkOptions {
	o.IgnoreErrors = ignore
	return o
}

// BulkGetPackages 批量获取多个包的信息
func (x *Repository) BulkGetPackages(ctx context.Context, gemNames []string, options ...*BulkOptions) []*BulkResult[*models.PackageInformation] {
	opts := getOptions(options)
	results := make([]*BulkResult[*models.PackageInformation], 0, len(gemNames))

	// 使用工作池模式处理
	workChan := make(chan string, len(gemNames))
	resultChan := make(chan *BulkResult[*models.PackageInformation], len(gemNames))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < opts.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for gemName := range workChan {
				pkg, err := x.GetPackage(ctx, gemName)
				resultChan <- &BulkResult[*models.PackageInformation]{
					Key:   gemName,
					Value: pkg,
					Error: err,
				}
			}
		}()
	}

	// 发送工作
	go func() {
		for _, gemName := range gemNames {
			select {
			case workChan <- gemName:
				// 成功添加到工作队列
			case <-ctx.Done():
				// 上下文被取消，停止添加
				break
			}
		}
		close(workChan)

		// 等待所有工作完成后关闭结果通道
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// BulkGetVersions 批量获取多个包的版本信息
func (x *Repository) BulkGetVersions(ctx context.Context, gemNames []string, options ...*BulkOptions) []*BulkResult[[]*models.Version] {
	opts := getOptions(options)
	results := make([]*BulkResult[[]*models.Version], 0, len(gemNames))

	// 使用工作池模式处理
	workChan := make(chan string, len(gemNames))
	resultChan := make(chan *BulkResult[[]*models.Version], len(gemNames))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < opts.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for gemName := range workChan {
				versions, err := x.GetGemVersions(ctx, gemName)
				resultChan <- &BulkResult[[]*models.Version]{
					Key:   gemName,
					Value: versions,
					Error: err,
				}
			}
		}()
	}

	// 发送工作
	go func() {
		for _, gemName := range gemNames {
			select {
			case workChan <- gemName:
				// 成功添加到工作队列
			case <-ctx.Done():
				// 上下文被取消，停止添加
				break
			}
		}
		close(workChan)

		// 等待所有工作完成后关闭结果通道
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// BulkGetDependencies 批量获取多个包的依赖信息
func (x *Repository) BulkGetDependencies(ctx context.Context, gemNames []string, options ...*BulkOptions) []*BulkResult[[]*models.DependencyInfo] {
	opts := getOptions(options)
	results := make([]*BulkResult[[]*models.DependencyInfo], 0, len(gemNames))

	// 使用工作池模式处理
	workChan := make(chan string, len(gemNames))
	resultChan := make(chan *BulkResult[[]*models.DependencyInfo], len(gemNames))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < opts.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for gemName := range workChan {
				deps, err := x.GetDependencies(ctx, gemName)
				resultChan <- &BulkResult[[]*models.DependencyInfo]{
					Key:   gemName,
					Value: deps,
					Error: err,
				}
			}
		}()
	}

	// 发送工作
	go func() {
		for _, gemName := range gemNames {
			select {
			case workChan <- gemName:
				// 成功添加到工作队列
			case <-ctx.Done():
				// 上下文被取消，停止添加
				break
			}
		}
		close(workChan)

		// 等待所有工作完成后关闭结果通道
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// BulkGetReverseDependencies 批量获取多个包的反向依赖
func (x *Repository) BulkGetReverseDependencies(ctx context.Context, gemNames []string, options ...*BulkOptions) []*BulkResult[[]string] {
	opts := getOptions(options)
	results := make([]*BulkResult[[]string], 0, len(gemNames))

	// 使用工作池模式处理
	workChan := make(chan string, len(gemNames))
	resultChan := make(chan *BulkResult[[]string], len(gemNames))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < opts.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for gemName := range workChan {
				deps, err := x.GetReverseDependencies(ctx, gemName)
				resultChan <- &BulkResult[[]string]{
					Key:   gemName,
					Value: deps,
					Error: err,
				}
			}
		}()
	}

	// 发送工作
	go func() {
		for _, gemName := range gemNames {
			select {
			case workChan <- gemName:
				// 成功添加到工作队列
			case <-ctx.Done():
				// 上下文被取消，停止添加
				break
			}
		}
		close(workChan)

		// 等待所有工作完成后关闭结果通道
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// 获取选项，如果未提供则使用默认值
func getOptions(options []*BulkOptions) *BulkOptions {
	if len(options) > 0 && options[0] != nil {
		return options[0]
	}
	return NewBulkOptions()
}
