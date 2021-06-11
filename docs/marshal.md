# Marshal & Unmarshal

阅读本章节之前您需要先行阅读
   
- [解析配置文件](./parse.md)
- [json.Marshal & json.Unmarshal][json_marshal]


如果您了解`json.Marshal`和`json.Unmarshal`那么本章节就简单的多（不了解？[点击了解][json_marshal]）。

由于配置文件的特殊性：`块`既可以包含`参数` 亦可以包含 `子块`。而在`json`或`yaml`中不存在这列情况。
所以我们分为两个部分讲解 Marshal:

- 简单模式：`块`存在参数就不存在`子块`, `Args`和`Body`不同时存在。
- 符合模式： 两者同时存在

## 名词解释

**基础类型**：int,[u]int[8|16|32|64],string,bool,time.Time类型以及对应的指针类型。

## 简单模式

加入我们现在有一个配置文件：

```nginx
user haiker;
address "Beijing China";
email ni@renzhen.la;
workyear 10;
Hobby basketball football "D D D" etc;
time "2021-06-11 12:13:14";
```

对应的的Golang结构体：
```nginx
type User struct {
    UserName string `ngx:"user"`
    Address string  `ngx:"address"`
    Email string    `ngx:"email"`
    Year int        `ngx:"workyear"`
    Hobby []string 
    Time time.Time  `ngx:"time,2016-01-02 15:04:05"`
} 
```

我们可以看到`Marshal`和`json.Marshal` 是一样的处理发方式。`略过。。。。。`，我们将把重点放在复杂模式上。

## 复杂模式

    复杂模式是针对一个特殊的类型做出的处理。他们包括了`slice`和`map`

### 切片Slice类型

#### 1、基础类型的切片

对于`[]int`,`[]string`,`[]bool`等基本类型的切片而言，我们只可以使用`块`的参数`args`赋值。

例如： 在简单模式中 `User.Hobby` 字段，我们可以使用下面方式方式赋予其值。

```nginx 
Hobby basketball football "D D D" etc;
``` 

也可以修改为多次指令赋予slice值。

```nginx
Hobby basketball; 
Hobby football;
Hobby "D D D" etc;
```
上面两种方式无论哪个对于最终结果是一致的。


#### 2、非基础型的切片

对于非基础类型的切片，就不可以像基础类型那样了两者均可使用。非基础类型只能采用多`块`配置方式。
例如：当配置 `Group`中包含`[]*User` 字段 我们就需要使用下面的配置。

```nginx
type Group struct {
    Users []*User `ngx:"users"`
}
```

配置如下
```nginx
users {
    user haiker;
    address "Beijing China";
    email ni@renzhen.la;
    workyear 10;
    Hobby basketball football "D D D" etc;
    time "2021-06-11 12:13:14";
}
users {
    user submitter;
    address "Beijing China";
    email demo@example.com;
    workyear 10;
    Hobby basketball;
    time "2021-06-11 12:13:14";
}
...
```

### map类型

Note: 配置文件转为map的类型，是有一定的规定，在map的key值类型只能为基础类型。

#### 1、对于value为基础类型

我们将直接采用`块`的`参数`进行值的赋予，且`参数1`为map的key, `参数2`为map的value。

```nginx
mapFiled mapKey mapValue;
mapFiled mapKey mapValue;
mapFiled mapKey mapValue;

mapFiled1 mapKey mapValue;
mapFiled1 mapKey mapValue;
mapFiled1 mapKey mapValue;
```

#### 2、对于value非基础类型

1、方式一：`块`的`参数1`作为key， `子块`作为value的配置。
```nginx
mapField key1 {
    valueField value;
    ...
}
mapField key2 {
    valueField value;
    ...
}
```
2、方式二：`块`的`子块`块名称作为key
```nginx 
mapField {
    key1 {
        valueField value;
        ...
    }
    key2 {
        valueField value;
        ...
    }
}
```

上面两种方式配置Unmarshal后对应的结果完全一致的。**如果采用了第二种方式配置map并且在子块上也添加了参数，将会报错。**
例如：
```nginx 
mapField {
    key1 arg0 { ## 这里不允许再次添加参数。。
        valueField value;
        ...
    }
    key2 {
        valueField value;
        ...
    }
}
```

## 重复覆盖

由于文件特性，很多时候可以指定相同的`块名称`的配置，配置文件本省并不存在错误，但是我们才用`Unmarshal`后对应的一个`基础类型`字段上的话，相当于对此赋值，最后出现的一次将生效，请注意这个默认选项。


## 自定义解析

如同`encoding/json`定义的`json.Unmarshaler`和`json.Marshaler`一样，
ngx也定义了`encoding.Marshaler`和`encoding.Unmarshalers` 方法, 定义如下：

```nginx
type Marshaler interface {
	MarshalNgx() (*config.Configuration, error)
}

type Unmarshaler interface {
	UnmarshalNgx(item *config.Configuration) error
}
```

如果你需要自定义解析方式实现该方法即可。例如：

```nginx

type User struct {
	UserName string `ngx:"name"`
	Address string  `ngx:"address"`
	Email string    `ngx:"email"`
	Year int        `ngx:"workyear"`
	Hobby []string
	Time time.Time  `ngx:"time,2016-01-02 15:04:05"`
}

func (user *User) MarshalNgx() (conf *config.Configuration,err error) {
	conf = &config.Configuration{
		Body: config.Directives{},
	}
	firstAndLastName := strings.SplitN(user.Name,".", 2)
	conf.Body = append(conf.Body,config.New("firstName",firstAndLastName[0]))
	conf.Body = append(conf.Body,config.New("lastName",firstAndLastName[1]))
	conf.Body = append(conf.Body,config.New("address",user.Address))
	...
	return
}

func (user *User) UnmarshalNgx(item *config.Configuration) (err error) {
	if len(item.Body) == 0 {
		return
	}
	
	if firstName := item.Body.Get("firstName"); firstName != nil {
		user.UserName = firstName.Args[0]	
	}
	if lastName := item.Body.Get("lastName"); lastName != nil {
		user.UserName = lastName.Args[0]
	}
	return
}
```

我们的`User`结构体，可以实现`encoding.Marshaler`和`encoding.Unmarshaler`方法去实现自定义的解析。

如果我们在配置中使用了第三方的库，此时无法让结构式实现接口，
因此系统同时也提供了`encoding.RegTypeHandler(v interface{}, handler encoding.TypeHandler)`方法指定特定第三方类型的解析方式。

`encoding.TypeHandler` 定义如下：

```nginx
TypeHandler interface {
    MarshalNgx(v interface{}) (*config.Configuration, error)
    UnmarshalNgx(item *config.Configuration) (v interface{}, err error)
}
```
更多玩法，赶快试用起来吧。


## 转JSON配置

本配置同时也提供了将配置转换为JSON配置[了解更多](./json.md)


>
> Note: 开源作者的英语能力有限，编写文档英文全靠翻译软件，难免存在歧义，
> 如果您可以为本软件编写英文文档请联系作者[Haiker](mailto:ni.renzhen.la)或提交PR，
> 感谢您对本工具的做出的贡献。
>

[json_marshal]: https://golang.org/pkg/encoding/json/ "JSON"
