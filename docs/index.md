# NGX

**NGX** 一个可以解析类似nginx.conf配置文件的Golang类库。

## 功能特性

- 解析类似 nginx.conf 的配置文件。
- 解析配置文件成Golang结构体，序列化Golang结构体到配置文件。
- 提供配置查询语言pql, 方便查询配置项查询.
- 配置文件动态化，根据参数生成不同的配置文件，模板化配置文件. 
- 解析后置hook处理操作，针对配置进行特出处理。
- 配置查询语言pql提供用户自定义扩展.
- 添加`NGXQ`命令行工具，为ngx配置文件的查询的命令行工具。
  使用方式类似于json的查询工具 [**JQ**][1] 和yaml的查询 [**YQ**][2]。
- 解析配置文件成json格式


## 使用指南（Starting Guide）

- [解析配置文件 parse configuration files](./parse.md)
- [系列化/反序列化 Marshal/Unmarshal](./marshal.md)
- [转换为JSON](./json.md)
- [配置查询语言 PQL](./pql.md)
- [后置操作 Post-hooks](./hooks.md)
- [命令工具 **NGX**](./cli.md)
