
## 使用指南

### 安装

```nginx
go get github.com/ihaiker/ngx/v2
```

### 实例

> Note: 为了方便，我们直接使用nginx的配置文件做讲解

加入我们已有`/etc/nginx/nginx.conf`配置文件，且内容如下：

```nginx
user nobody;
worker_processes auto;
events  {
    worker_connections 1024;
}
http  {
    include mime.types;
    default_type application/octet-stream;
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
					'$status $body_bytes_sent "$http_referer" '
					'"$http_user_agent" "$http_x_forwarded_for"';
    # access_log  /var/log/nginx/access.log  main;
    sendfile on;
    # tcp_nopush     on;
    keepalive_timeout 65;
    gzip on;
    include conf.d/*.conf;
}
```

解析Golang代码如下： 

```nginx

package main

import (
	"fmt"
	"github.com/ihaiker/ngx/v2"
)

func main() {
	conf, err := ngx.Parse("/etc/nginx/nginx.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("--- configuration body ---")
	for i, item := range conf.Body {
		fmt.Println(i, ": ", item.Name)
	}

	fmt.Println("--- http body ---")
	for i, item := range conf.Body[3].Body {
		fmt.Println(i, ": ", item.Name)
	}
}

```
run it, the following will be print:

```shell 
--- configuration body ---
0 :  user
1 :  worker_processes
2 :  events
3 :  http
--- http body ---
0 :  include
1 :  default_type
2 :  log_format
3 :  #
4 :  sendfile
5 :  #
6 :  keepalive_timeout
7 :  gzip
8 :  include
include args:  mime.types
```

### 结构说明

首先，我们来分析一下配置文件的结构。结构有以下规则：

```nginx 
name [arg0 arg1 arg3 ...];
name [arg0 arg1 ... argN] {
    name [arg0 arg1 ...];
    name [arg0 arg1 ...];
}
```
由上结构可分析出，组成他的最小化单元是块，这个块由`块名称` **name**，`块参数` **arg0 arg1 arg2**(可选的), `子块`（可选）组成，
且`子块`结构也是按照块的结构组成，依次类推。

因此，我们将配置文件解析为 Golang 结构为:

```nginx
type Directive struct {
    Line    int        `json:"line"`
    Virtual Virtual    `json:"virtual,omitempty"`
    Name    string     `json:"name"`
    Args    []string   `json:"args,omitempty"`
    Body    Directives `json:"body,omitempty"`
}
```

在 Golang struct `Directive` 中包含 `Line`,`Name`,`Args`,`Body` 字段.
`Line` 字段保存了当前块出现在文件的行号， 
`Name` 字段存放块的名称, 
`Args` 字段为块的参数,
`Body` 字段为子块.
`Virtual` 字段有些特出，为[Post-hook](./hooks.md)方法准备，用于解决include这样的特出指令。

调用解析后程序会生成`config.Configuration`结构体，保存配置信息。
```nginx
type (
    Directives    []*Directive
    Configuration struct {
        Source string     `json:"source"` //配置文件来源：files://, stdin, parse: hql
        Body   Directives `json:"body"`
    }
)
```

但是，我们也可以发现，到现在虽然我们已经将文件解析了，但是他的使用上来看并不是很优雅，
我们查找一个字段的时候需要多次循序，例如：`conf.Body[2].Body[1].Args[0]` 因此我们编写了方便的处理方式：

- [Marshal/Unmarshal指南](./marshal.md) 整合到自定义结构体中。
- [PQL查询语言](./pql.md)  配置查询语言，简单快捷查询

>
> Note: 开源作者的英语能力有限，编写文档英文全靠翻译软件，难免存在歧义，
> 如果您可以为本软件编写英文文档请联系作者[Haiker](mailto:ni.renzhen.la)或提交PR，
> 感谢您对本工具的做出的贡献。
>
