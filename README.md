# NGX

## 简介

**NGX** 一个可以解析类似nginx配置文件的Golang类库。

## 特性

- 解析类似 nginx.conf 的配置文件。
- 解析配置文件成Golang结构体，序列化Golang结构体到配置文件。
- 提供配置查询语言pql, 方便查询配置项查询.
- 配置文件动态化，根据参数生成不同的配置文件，模板化配置文件. 
- 解析后置hook处理操作，针对配置进行特出处理。
- 配置查询语言pql提供用户自定义扩展.
- 添加`NGXQ`命令行工具，为ngx配置文件的查询的命令行工具。
  使用方式类似于json的查询工具 [**JQ**][1] 和yaml的查询 [**YQ**][2]。
- 解析配置文件成json格式

## 使用帮助 

查看本项目 [wiki](https://ihaiker.github.io/ngx) 获取使用帮助.

文档地址：[https://ihaiker.github.io/ngx](https://ihaiker.github.io/ngx)


## Q&A

欢迎任何问题或问题，请提交它们 [github issuse](https:github.comihaikerngxissues) :)

## 使用软件列表

- [aginx](https://github.com/ihaiker/aginx) 为`nginx`添加管理restful api和控制台
- [vik8s](https://github.com/ihaiker/vik8s) `k8s` 集群安装工具


[1]: https://stedolan.github.io/jq/ "JQ"
[2]: https://mikefarah.gitbook.io/yq/ "YQ"
