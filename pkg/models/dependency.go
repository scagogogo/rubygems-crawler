package models

// DependencyInfo 用于/api/v1/dependencies接口
// 参考: https://guides.rubygems.org/rubygems-org-api-v2/#dependencies
type DependencyInfo struct {
	// 包名
	Name string `json:"name"`

	// 依赖的包名
	DependentName string `json:"dependent_name"`

	// 版本要求，例如: ">= 1.0.0"
	Requirements string `json:"requirements"`

	// 依赖类型，常见值: "runtime", "development"
	DependentType string `json:"dependent_type"`
}
