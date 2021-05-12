package config

type Options struct {
	Delimiter     bool //是否允许分割符号":"。例如： server_name: nginx; 方式是否允许使用。
	RemoveQuote   bool //去除括号, 如果为true, 例如： name "aginx.io"; 这里的括号是否去除
	RemoveCommits bool //去除注解内容，如果为true将不再保留注释内容
	/**合并Include内容，如果使用后竞不在含有 include 内容而是直接吧 include内容加载过来*/
	MergeInclude bool
}

var def = Default()

func Default() *Options {
	return &Options{
		Delimiter:     true,
		RemoveQuote:   false,
		RemoveCommits: false,
		MergeInclude:  false,
	}
}
