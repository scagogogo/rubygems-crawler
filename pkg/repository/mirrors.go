package repository

// ------------------------------------------------- --------------------------------------------------------------------

const ServerURLRubyChina = "https://gems.ruby-china.com"

// NewRubyChinaRepository 使用Ruby中国的镜像仓库，国内用户推荐使用这个镜像源
func NewRubyChinaRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(ServerURLRubyChina))
}

// ------------------------------------------------- --------------------------------------------------------------------

// TODO 2023-5-29 00:37:26 清华源测试不通过
const ServerURLTSingHua = "https://mirrors.tuna.tsinghua.edu.cn/rubygems"

// NewTSingHuaRepository 使用清华大学的镜像仓库
func NewTSingHuaRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(ServerURLTSingHua))
}

// ------------------------------------------------- --------------------------------------------------------------------


// 阿里云
// https://mirrors.aliyun.com/rubygems/
