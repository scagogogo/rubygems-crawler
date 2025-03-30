package repository

// ------------------------------------------------- --------------------------------------------------------------------

const ServerURLRubyChina = "https://gems.ruby-china.com"

// NewRubyChinaRepository 使用Ruby中国的镜像仓库，国内用户推荐使用这个镜像源
func NewRubyChinaRepository() Repository {
	return NewRepository(NewOptions().SetServerURL(ServerURLRubyChina))
}

// ------------------------------------------------- --------------------------------------------------------------------

// 之前的清华源API路径有误，更新为正确的URL
// 原问题：清华源使用了不同的API路径结构
const ServerURLTSingHua = "https://mirrors.tuna.tsinghua.edu.cn/rubygems/api"

// NewTSingHuaRepository 使用清华大学的镜像仓库
func NewTSingHuaRepository() Repository {
	return NewRepository(NewOptions().SetServerURL(ServerURLTSingHua))
}

// ------------------------------------------------- --------------------------------------------------------------------

// 阿里云
// https://mirrors.aliyun.com/rubygems/

// ------------------------------------------------- --------------------------------------------------------------------

const ServerURLAliYun = "https://mirrors.aliyun.com/rubygems"

// NewAliYunRepository 使用阿里云的镜像仓库
func NewAliYunRepository() Repository {
	return NewRepository(NewOptions().SetServerURL(ServerURLAliYun))
}
