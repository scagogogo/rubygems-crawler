package repository

// DefaultServerURL 默认的仓库地址，直接连接到官方仓库
const DefaultServerURL = "https://rubygems.org"

type Options struct {

	// 仓库的服务器地址
	ServerURL string

	// 向仓库发送请求时使用代理
	Proxy string
}

func NewOptions() *Options {
	return &Options{
		ServerURL: DefaultServerURL,
		Proxy:     "",
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
