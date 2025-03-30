package repository

import (
	"context"
	"sync"

	"github.com/scagogogo/rubygems-crawler/pkg/models"
)

// BulkResult 表示批量操作的单个结果
// 它包含请求的键（如包名）、返回的值和可能的错误
type BulkResult[T any] struct {
	Key   string // 请求的键，通常是gem包名
	Value T      // 操作的结果值
	Error error  // 操作过程中可能发生的错误
}

// BulkOptions 定义批量操作的配置选项
type BulkOptions struct {
	// MaxConcurrency 定义最大并发请求数量
	// 设置合理的值可以避免对API服务器造成过大压力
	// 默认值为10
	MaxConcurrency int

	// ContinueOnError 决定在遇到错误时是否继续处理其他请求
	// 如果为true，即使某些请求失败，仍会处理所有请求
	// 如果为false，遇到第一个错误时会立即停止处理
	// 默认为true
	ContinueOnError bool
}

// NewBulkOptions 创建具有默认值的批量操作选项
// 默认配置：最大并发数10，遇到错误时继续执行
func NewBulkOptions() *BulkOptions {
	return &BulkOptions{
		MaxConcurrency:  10,
		ContinueOnError: true,
	}
}

// WithMaxConcurrency 设置最大并发请求数
// 返回选项对象自身，支持链式调用
func (o *BulkOptions) WithMaxConcurrency(maxConcurrency int) *BulkOptions {
	if maxConcurrency > 0 {
		o.MaxConcurrency = maxConcurrency
	}
	return o
}

// WithContinueOnError 设置遇到错误时是否继续
// 返回选项对象自身，支持链式调用
func (o *BulkOptions) WithContinueOnError(continueOnError bool) *BulkOptions {
	o.ContinueOnError = continueOnError
	return o
}

// BulkGetPackages 批量获取多个包的信息
// 并发执行GetPackage请求，提高大规模数据获取效率
// 参数:
//   - ctx: 上下文，用于控制请求超时和取消
//   - gemNames: 要获取的包名列表
//   - options: 批量操作选项，控制并发数等
//
// 返回:
//   - 包含每个包请求结果的切片，顺序与输入包名相同
func (r *RepositoryImpl) BulkGetPackages(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[*models.PackageInformation] {
	if options == nil {
		options = NewBulkOptions()
	}

	results := make([]*BulkResult[*models.PackageInformation], len(gemNames))

	// 创建工作池
	worker := func(wg *sync.WaitGroup, jobs <-chan int, results []*BulkResult[*models.PackageInformation]) {
		defer wg.Done()

		for i := range jobs {
			select {
			case <-ctx.Done():
				// 上下文被取消，停止处理
				results[i] = &BulkResult[*models.PackageInformation]{
					Key:   gemNames[i],
					Error: ctx.Err(),
				}
				return
			default:
				// 获取包信息
				pkg, err := r.GetPackage(ctx, gemNames[i])
				results[i] = &BulkResult[*models.PackageInformation]{
					Key:   gemNames[i],
					Value: pkg,
					Error: err,
				}

				// 如果设置了遇到错误停止，并且发生了错误
				if !options.ContinueOnError && err != nil {
					return
				}
			}
		}
	}

	// 运行工作池
	runWorkerPool(options.MaxConcurrency, len(gemNames), results, worker)

	return results
}

// BulkGetVersions 批量获取多个包的版本信息
// 并发执行GetGemVersions请求，提高大规模数据获取效率
// 参数:
//   - ctx: 上下文，用于控制请求超时和取消
//   - gemNames: 要获取的包名列表
//   - options: 批量操作选项，控制并发数等
//
// 返回:
//   - 包含每个包版本请求结果的切片，顺序与输入包名相同
func (r *RepositoryImpl) BulkGetVersions(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]*models.Version] {
	if options == nil {
		options = NewBulkOptions()
	}

	results := make([]*BulkResult[[]*models.Version], len(gemNames))

	// 创建工作池
	worker := func(wg *sync.WaitGroup, jobs <-chan int, results []*BulkResult[[]*models.Version]) {
		defer wg.Done()

		for i := range jobs {
			select {
			case <-ctx.Done():
				// 上下文被取消，停止处理
				results[i] = &BulkResult[[]*models.Version]{
					Key:   gemNames[i],
					Error: ctx.Err(),
				}
				return
			default:
				// 获取版本信息
				versions, err := r.GetGemVersions(ctx, gemNames[i])
				results[i] = &BulkResult[[]*models.Version]{
					Key:   gemNames[i],
					Value: versions,
					Error: err,
				}

				// 如果设置了遇到错误停止，并且发生了错误
				if !options.ContinueOnError && err != nil {
					return
				}
			}
		}
	}

	// 运行工作池
	runWorkerPool(options.MaxConcurrency, len(gemNames), results, worker)

	return results
}

// BulkGetDependencies 批量获取多个包的依赖信息
// 并发执行GetDependencies请求，提高大规模数据获取效率
// 参数:
//   - ctx: 上下文，用于控制请求超时和取消
//   - gemNames: 要获取的包名列表
//   - options: 批量操作选项，控制并发数等
//
// 返回:
//   - 包含每个包依赖请求结果的切片，顺序与输入包名相同
func (r *RepositoryImpl) BulkGetDependencies(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]*models.DependencyInfo] {
	if options == nil {
		options = NewBulkOptions()
	}

	results := make([]*BulkResult[[]*models.DependencyInfo], len(gemNames))

	// 创建工作池
	worker := func(wg *sync.WaitGroup, jobs <-chan int, results []*BulkResult[[]*models.DependencyInfo]) {
		defer wg.Done()

		for i := range jobs {
			select {
			case <-ctx.Done():
				// 上下文被取消，停止处理
				results[i] = &BulkResult[[]*models.DependencyInfo]{
					Key:   gemNames[i],
					Error: ctx.Err(),
				}
				return
			default:
				// 获取依赖信息
				deps, err := r.GetDependencies(ctx, gemNames[i])
				results[i] = &BulkResult[[]*models.DependencyInfo]{
					Key:   gemNames[i],
					Value: deps,
					Error: err,
				}

				// 如果设置了遇到错误停止，并且发生了错误
				if !options.ContinueOnError && err != nil {
					return
				}
			}
		}
	}

	// 运行工作池
	runWorkerPool(options.MaxConcurrency, len(gemNames), results, worker)

	return results
}

// BulkGetReverseDependencies 批量获取多个包的反向依赖信息
// 并发执行GetReverseDependencies请求，提高大规模数据获取效率
// 参数:
//   - ctx: 上下文，用于控制请求超时和取消
//   - gemNames: 要获取的包名列表
//   - options: 批量操作选项，控制并发数等
//
// 返回:
//   - 包含每个包反向依赖请求结果的切片，顺序与输入包名相同
func (r *RepositoryImpl) BulkGetReverseDependencies(ctx context.Context, gemNames []string, options *BulkOptions) []*BulkResult[[]string] {
	if options == nil {
		options = NewBulkOptions()
	}

	results := make([]*BulkResult[[]string], len(gemNames))

	// 创建工作池
	worker := func(wg *sync.WaitGroup, jobs <-chan int, results []*BulkResult[[]string]) {
		defer wg.Done()

		for i := range jobs {
			select {
			case <-ctx.Done():
				// 上下文被取消，停止处理
				results[i] = &BulkResult[[]string]{
					Key:   gemNames[i],
					Error: ctx.Err(),
				}
				return
			default:
				// 获取反向依赖信息
				deps, err := r.GetReverseDependencies(ctx, gemNames[i])
				results[i] = &BulkResult[[]string]{
					Key:   gemNames[i],
					Value: deps,
					Error: err,
				}

				// 如果设置了遇到错误停止，并且发生了错误
				if !options.ContinueOnError && err != nil {
					return
				}
			}
		}
	}

	// 运行工作池
	runWorkerPool(options.MaxConcurrency, len(gemNames), results, worker)

	return results
}

// runWorkerPool 是一个通用的工作池实现，用于并发处理任务
// 参数:
//   - numWorkers: 工作协程数量
//   - numJobs: 总任务数量
//   - results: 存储结果的切片
//   - workerFunc: 工作函数，定义了每个工作协程的行为
func runWorkerPool[T any](numWorkers, numJobs int, results []*BulkResult[T], workerFunc func(*sync.WaitGroup, <-chan int, []*BulkResult[T])) {
	// 确保工作协程数量不超过任务数量
	if numWorkers > numJobs {
		numWorkers = numJobs
	}

	// 创建工作组和任务通道
	var wg sync.WaitGroup
	jobs := make(chan int, numJobs)

	// 启动工作协程
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go workerFunc(&wg, jobs, results)
	}

	// 分发任务
	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// 等待所有工作协程完成
	wg.Wait()
}
