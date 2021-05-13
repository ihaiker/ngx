package config

type Options struct {
	Delimiter     bool //是否允许分割符号":"。例如： server_name: nginx; 方式是否允许使用。
	RemoveQuote   bool //去除括号, 如果为true, 例如： name "aginx.io"; 这里的括号是否去除
	RemoveCommits bool //去除注解内容，如果为true将不再保留注释内容
}

var def = Default()

func Default() *Options {
	return &Options{
		Delimiter:     true,
		RemoveQuote:   false,
		RemoveCommits: false,
	}
}

//for nginx config file options
func Nginx() *Options {
	return &Options{
		Delimiter:     false,
		RemoveQuote:   false,
		RemoveCommits: true,
	}
}

//Encoding 为 encoding的options
func Encoding() *Options {
	return &Options{
		Delimiter:     true,
		RemoveQuote:   true,
		RemoveCommits: true,
	}
}
