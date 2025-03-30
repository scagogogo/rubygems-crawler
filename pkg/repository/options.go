package repository

// DefaultServerURL 默认的仓库地址，直接连接到官方仓库
const DefaultServerURL = "https://rubygems.org"

type Options struct {

	// 仓库的服务器地址
	ServerURL string

	// 向仓库发送请求时使用代理
	Proxy string

	// 用于API认证的Token
	// 参考: https://guides.rubygems.org/rubygems-org-api-v2/#rate-limits
	Token string

	// 请求重试选项
	RetryOptions *RetryOptions
}

func NewOptions() *Options {
	return &Options{
		ServerURL:    DefaultServerURL,
		Proxy:        "",
		Token:        "",
		RetryOptions: NewDefaultRetryOptions(),
	}
}

func (x *Options) SetServerURL(serverUrl string) *Options {
	x.ServerURL = serverUrl
	return x
}

func (x *Options) SetProxy(proxy string) *Options {
	x.Proxy = proxy
	return x
}

func (x *Options) SetToken(token string) *Options {
	x.Token = token
	return x
}

func (x *Options) SetRetryOptions(retryOptions *RetryOptions) *Options {
	x.RetryOptions = retryOptions
	return x
}

// DisableRetry 禁用重试功能
func (x *Options) DisableRetry() *Options {
	x.RetryOptions = nil
	return x
}
