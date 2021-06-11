# PQL 配置查询语言

阅读本章节内容之前您必须先行了解：[配置文件解析](./parse.md)

如果您阅读了之前的章节，您会发现配置在使用上其实并不是特别简单和明了，
我们提供了[Marshal/Unmarshal](./marshal.md)方式让使用更简单，
同时也提供了配置搜索功能，给他在插上一双翅膀。

## 实例Demo
为了让本章节所有的讲解可以更容易理解，我们使用一个`nginx.conf`来说明：

```nginx
user  nginx;
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

    include hosts.d/*.conf;

    upstream t1 {
        server 127.0.0.1:8001;
    }
    upstream t2 {
        server 127.0.0.1:8002;
    }
    upstream t3 {
        server 127.0.0.1:8002;
    }
    server  {
        listen 80;
        server_name _;
        location / {
            root /Users/haiker/Documents/bootstramp/coreui/dist;
            index index.html index.htm;
        }
        location /health {
            return 200 'ok';
        }
    }

    server  {
        listen 80;
        server_name aginx.x.do;
        location / {
            proxy_pass http://t1;
            gzip on;
        }
    }

    server  {
        listen 80;
        server_name test.renzhen.la;
        location / {
            proxy_pass http://t2;
            proxy_set_header Host $host;
        }
    }
}
```

先来个demo来看看：

```nginx 
package main

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/query"
)

func main() {
	conf, err := config.Parse("./nginx.conf")
	if err != nil {
		panic(err)
	}
    // ----- 执行查询 ------
	items, err := query.Select(conf, ".http.upstream")
	if err != nil {
		panic(err)
	}
	fmt.Println("包含 upstream：", len(items))

	upstreamT1 := items[0]
	fmt.Println("upstream t1 name: ",upstreamT1.Args[0])
	fmt.Println("upstream t1 server value:",upstreamT1.Body.Get("server").Args[0])

	upstreamT2 := items[1]
	fmt.Println("upstream t2 name: ",upstreamT2.Args[0])
	fmt.Println("upstream t2 server value:",upstreamT2.Body.Get("server").Args[0])

	upstreamT3 := items[2]
	fmt.Println("upstream t3 name: ",upstreamT3.Args[0])
	fmt.Println("upstream t3 server value:",upstreamT3.Body.Get("server").Args[0])
}
```

## 语法介绍

看了这个demo你你是不要有点对于pql感觉有熟悉的感觉呢？没错，本工具在设计上借鉴了json查询工具[jq][JQ]和yaml查询工具[yq][YQ]。并且语法上基本一致。


下面我们就来详细说明PQL语法。

### ` . ` 指令

直接查询配置文件当前节点。

### `.name`或`.name.subname` 指令 

查询`name`或者`name`底下`subname`配置项，因为`ngx`配置可以允许配置项名称重复，所以无论如何查询都返回的内容都是数组。

例如：
```shell 
.http.server.server_name
``` 
可以查询nginx配置文件中所有的server_name，输入如下：

```nginx
server_name _;
server_name aginx.x.do;
server_name test.renzhen.la;
```

### `.!name`，`.@name`，`.^name`，`.$name` 指令

|表述符|描述|
|---|---|
|!|配置名称不等于name的配置项|
|@|配置名称包含name的配置项|
|^|配置名称以name开头的配置项|
|$|配置名称以name结尾的配置项|

### `./Regex/` 指令

查找配置项名称符合**Regex**正则的配置项。

例如：
```shell
.http./(server|upstream)/
```

查询**http**下面所有的**server**和**upstream**


### `.name[0]` 指令

查询配置名称为`name`的配置项的**第一个**，如果并未能查询到相应配置将会报错。

为了查询方便软件也提供了以下方式：
- `.name[start:end]`查询某个区间的配置。
- `.name[start:]`查询从某位置开始区间配置。
- `.name[:end]`查询开始到某个位置的区间配置。
- `.name[-1]`查询最后一个配置。

例如：
```shell
.http.server[1].server_name
```
输出：
```nginx
server_name aginx.x.do;
```

### `.name("value")` 指令 

查询配置中，包含value的配置项。

例如：
```shell
.http.server.server_name("_")
```
输出：
```nginx
server_name _;
```

此处的value也同时支持名称查询的`!`,`@`,`^`,`$`以及`/value/`正则五种查询方式，且具有相同意思。
例如：
```shell
.http.server.server_name($"example.com")
```
查询nginx配置中使用了**example.com**一级域名的server配置

```shell
.http.server.server_name(^"test.")
```
查询server配置中，test开始的域名配置


### " | " 管道处理

根据上次结果继续处理。
`.http | .server_name` 和 `.http.server_name`等价。但是只是用例好像没有什么多大用处，不过此处大多是为了查询方法准备的。
例如：
```shell 
.http.server | select(.server_name, "equal", "_")
```
查询**server_name**为**_**的server配置项。

## 系统查询方法

系统提供了为数不多的方法，详见下列表：

- [select](./functions/select.md) 
- [ifelse,not,and,or](./functions/grammar.md)
- [args](./functions/args.md)
- [index](./functions/index.md)
- [length](./functions/length.md)
- [args_length](functions/args_length.md)

注意：在上面的例子中估计你已经发了一个问题。就是select方法并不会项JQ或者YQ那样直接使用**select(.server_name equal _)** 这样的语法，
主要原因是兼容所有定义的方法具有相同的语法。（这个地方需要优化...下个版本解决吧）

## 自定义查询方法

为了使用自定义查询方法，你需要使用`query.SelectsWith(*config.Configuration, *methods.FunctionManager,string)` 方法，并指定查询方法管理器`*methods.FunctionManager`。
查询方法管理器可以注册任何方法，只要符合一下条件即可。

- 输入参数只能为：`error`, `string`, `int`, `config.Directive`, `bool`, `interface {}` 类型，或者对应的切片类型。
- 输出参数最多可以有两个，且第二个必须是error类型，第一个要符合和时输入参数一样的类型。



[JQ]: https://stedolan.github.io/jq/ "JQ"
[YQ]: https://mikefarah.gitbook.io/yq/ "YQ"
