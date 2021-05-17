
# 语法介绍

[!^$@]Ident[(arg1[\&\&|\|\|()])]+

[!^$@]Ident[(\.arg)]+

1、http、!http、^http、$http、@http、'@http',"!http"
2、http(name), http(name1 || name2)、http(name1 && name2)

ngx http server.server_name(baidu.com) 

####  定位参数详解

看完上面几个简单的示例，您一定对定位参数`q`有些疑惑，他到底值如何工作的。此小结我们将来介绍一下他的工作原理。

在解释参数`q`之前我能先来了解另外一个内容**NGINX配置**。

在对nginx配置分析后此软件将解析文件为：一个树形结构。并且这个树形结果的每个节点都是相同的结构。

分析的结果：

```nginx
directive [arg0 arg1 arg3 ... argN];
directive [arg0 arg1 ... argN ;] [{
      directive arg0 arg1 ... argN;
      directive arg0 arg1 ... argN;
}]
```

| 节点内容    | 说明                                                        |
| ----------- | ----------------------------------------------------------- |
| directive   | 指令，例如: http,server,server_name,tcp_nopush              |
| arg0...argN | 指令参数：例如：**listen 80**中的80，**sendfile on** 中的on |
| [{...}]     | 指令的下级指令。                                            |

**定位参数** 语法：

```javascript
[comparison]directive(comparison'arg1'[ operator comparison'arg2']+)[.[comparison]directive(comparison'arg1'[ operator comparison'arg2']+)]+
```

语法内容：

| 语法指令    | 可用值                         | 语法内容说明                                                 |
| ----------- | ------------------------------ | ------------------------------------------------------------ |
| comparison  | 指令或者参数的比较方式。       | 比较方式供五种<br />空：相等<br />! : 不相等<br />@: 包含<br />^: hasPrefix<br />$: hasSubffix |
| directive   | 指令名称                       | http,server等nginx的指令                                     |
| arg0...argN | 参数名称                       | nginx指令的参数                                              |
| operator    | 参数匹配的结果级的合并判断方式 | ”&“ 并且  ，“\|” 或者                                        |

估计您看到这里这个语法会有些晦涩难懂，那我们来拆解一下这个语法。

上面的语法其实是 **query.subQuery.subQuery1** 这个query就是语法的一个最小单元：

```json
[comparison]directive(comparison'arg1'[ operator comparison'arg2']+)
```

那么看懂了这个语法，接下来给出一些简单的实例让你理解更佳深入一些。



**1、查询 http下面所有server 配置 参数如下，**

```json
q=http
q=server
```

**2、查询http下面 server_name 为 portainer.aginx.io的server** 

```
q=http
q=server.server_name('portainer.aginx.io')
```

**3、查询http下面 server_name 为 portainer.aginx.io 并且 使用了ssl 的 server** 

```
q=http
q=server.[server_name('portainer.aginx.io') & listen('80')]
```

**3、查询 http server下server_name 以www开头的server** 

```json
q=http
q=server.server_name(^'www')
```

**4、查询 api.aginx.io 配置内容**

```json
q=http
q=include
q=file('hosts.d/api.conf')
q=server.server_name('hosts.d/api.conf')
```

**5、查询所有的server**

```json
q=http
q=server
```

如果使用include使用下面的。

```json
q=http
q=include
q=*
q=server
```
