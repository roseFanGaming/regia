## 快速开始

```golang
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
	engine.GET("/", func(ctx *regia.Context) {
		ctx.Response.String("hello world")
	})
	engine.Run(":8000")
}
```





### Engine

`Engine`是一个运行时核心组件的集合。

* `Router`：负责注册和匹配路由

* `HtmlRender`：负责渲染`HTML`页面

* `JsonSerializer`：负责`json`的序列化和反序列化

* `XmlSerializer`：负责`XML`的序列化和反序列化

* `FileStorage`：负责文件存储

* `Abort`：当调用`Context.Abort`的时候会调用该接口的`Exit`方法，默认什么都不做。

* `NotFoundHandle`：当路由匹配不到的时候会调用该方法。

* `Interceptors`：全局拦截器的集合，无论路由是否匹配上，拦截器都会被执行。

* `Starters`：`Starter`的集合，所有的`Starter`会在调用`Engine.Run`的时候执行。

* `Warehouse`：用来往`Engine`里存储信息

* `MultipartFormMaxMemory` : 设置`multipart form max size`

  

#### New

```go
func New() *Engine
```

`New`方法用来创建一个完善的`Engine`对象，它全部使用系统自带的组件。

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.New()
	engine.Run(":8000")
}
```



#### Default

```go
func Default() *Engine
```

`Default`方法在`New`的基础上加上了一些简单的功能，参考`Default`的实现来定制自己的`Engine`

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
	engine.Run(":8000")
}
```



#### Static

```go
func (e *Engine) Static(url, dir string, group ...HandleFunc)
```

`Static`方法提供静态文件服务，如：**images**、 **Css**、 **JavaScript**

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
	engine.Static("/static/", "static")
	engine.Run(":8000")
}

// http://localhost:8000/static/a.png
```



#### AddInterceptors

```go
func (e *Engine) AddInterceptors(interceptors ...HandleFunc)
```

`AddInterceptors`方法用来添加全局请求拦截器，无论路由是否匹配上，添加的拦截器都会被执行, 且拦截器会优先中间件和`Handle`执行



#### AddStarter

```go
func (e *Engine) AddStarter(starters ...Starter)
```

`AddStarter`方法用来添加`Starter`，`Starter`会在项目运行时启动，并且只会运行一次。



#### Run

```go
func (e *Engine) Run(addr string) error
```

`Run`方法用来启动当前的服务



### Branch

路由分支



#### NewBranch

```go
func NewBranch() *Branch
```

`NewBranch`用来创建一个`Branch`对象



#### Use

```go
func (b *Branch) Use(group ...HandleFunc)
```

注册当前分支的中间件, 注册的中间件只会在当前分支被匹配到的时候执行。



#### SetPrefix

```go
func (b *Branch) SetPrefix(path string)
```

设置当前的分支的前缀

```go
package main

import (
	"github.com/eatMoreApple/regia"
	"net/http"
)

func main() {
	engine := regia.Default()
	b := regia.NewBranch()
	b.SetPrefix("/user")
	b.Handle(http.MethodGet, "/login", func(ctx *regia.Context) {
		ctx.Response.String("login page")
	})
	engine.Include("", b)
	engine.Run(":8000")
}

// => http://localhost:8000/user/login
```



#### Handle

```go
func (b *Branch) Handle(method, path string, group ...HandleFunc)
```

注册自定义方法

```go
b := regia.NewBranch()
b.Handle(http.MethodGet, "/login", func(ctx *regia.Context) {
		ctx.Response.String("login page")
})
```



#### Include

嵌套另一个`Branch`

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
	b := regia.NewBranch()
	b.GET("/login", func(ctx *regia.Context) {
		ctx.Response.String("login page")
	})
	engine.Include("/user", b)
	engine.Run(":8000")
}

// => http://localhost:8000/user/login
```



#### Bind

```go
func (b *Branch) Bind(path string, v interface{}, mappings ...map[string]string) 
```

`Bind`方法用来注册`Struct`路由，通过`mappings`的请求方法和`handle`名称的映射来完成请求的分发。

**注**：`mappings`注册的方法必须是有效的方法`func(*regia.Context)`，否则不会被注册

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

type User struct{}

func (u *User) Login(ctx *regia.Context) {
	ctx.Response.String("login page")
}

func (u *User) Register() {}

func main() {
	engine := regia.Default()
	b := regia.NewBranch()
	b.Bind("/login", &User{}, map[string]string{"Login": "post", "Register": "post"})
	engine.Include("/user", b)
	engine.Run(":8000")
}

// POST /user/login
```



#### BindMethod

```go
func (b *Branch) BindMethod(path string, v interface{}, mappings ...map[string]string)
```

`BindMethod`的方法签名和`Bind一致`,但是`BindMethod`默认提供了一组`httpMethods`的方法映射, 详情请看`regia.HttpRequestMethodMapping`,可以在此基础上自定义方法的映射。

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

type User struct{}

func (u *User) Post(ctx *regia.Context) {
	ctx.Response.String("login page")
}

func (u *User) Detail(ctx *regia.Context) {
	ctx.Response.String("detail page")
}

func main() {
	engine := regia.Default()
	b := regia.NewBranch()
	b.BindMethod("/login", &User{}, map[string]string{"Detail": "get"})
	engine.Include("/user", b)
	engine.Run(":8000")
}

// GET => /user/login
// POST => /user/login
```



### Context

`Context`是一个用来控制请求、响应、连接`Engine`的上下文管理。



#### Next

```go
func (c *Context) Next()
```

在当前的`handle`中直接调用下一层的`handle`



#### Abort

```go
func (c *Context) Abort() 
```

结束并退出当前及下层的`handle`处理流程

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

func abort(ctx *regia.Context) {
	ctx.Abort()
	ctx.Next() // 不会被执行了
}

func main() {
	engine := regia.Default()
	engine.GET("/", abort, func(ctx *regia.Context) {
		ctx.Response.String("hello world")
	})
	engine.Run(":8000")
}
```

重置`Engine.Abort`来修改默认的`Abort`行为



#### AbortWith

```go
func (c *Context) AbortWith(exit Exit)
```

实现`Exit`接口，在结束并退出`Handle`处理流程前后执行`Exit.Exit`

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

type MyExit struct{}

func (m MyExit) Exit(ctx *regia.Context) {
	ctx.Response.String("forbidden")
}

func main() {
	engine := regia.Default()
	engine.GET("/", func(ctx *regia.Context) {
		ctx.AbortWith(MyExit{})
		ctx.Response.String("hello world")
	})
	engine.Run(":8000")
}
```



#### SaveUploadFile

```go
func (c *Context) SaveUploadFile(filer *File, path string) error
```

保存上传的文件，这个方法将会调用`Engine`的`FileStorage.Save`。

`Engine`默认的`FileStorage`是将文件存储到本地，可以重置这个接口来实现自定义文件保存。

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
	engine.POST("/", func(ctx *regia.Context) {
		files, _ := ctx.Request.Files()
		file, _ := files.Get("file")
		ctx.SaveUploadFile(file, "your file name")
		ctx.Response.String("upload success")
	})
	engine.Run(":8000")
}
```



### Request



#### Query

```go
func (r *Request) Query() URLValue
```

获取当前的URL参数

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
	engine.GET("/", func(ctx *regia.Context) {
		name := ctx.Request.Query().Get("name")
		result := name.String("eatMoreApple") // can set default value
		ctx.Response.String(result)
	})
	engine.Run(":8000")
}

// => http://localhost:8000/?name=ivy 
// ivy

// => http://localhost:8000/
// eatMoreApple
```



#### Form

```go
func (r *Request) Form() URLValue
```

`Form`方法可以获取请求方式为`POST`、`PUT`、`PATCH`方法的`form data`

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
	engine.POST("/", func(ctx *regia.Context) {
		name := ctx.Request.Form().Get("name")
		result := name.String("eatMoreApple") // set default value
		ctx.Response.String(result)
	})
	engine.Run(":8000")
}
```



#### Files

```go
func (r *Request) Files() (Files, error)
```

获取所有上传的文件

```go
package main

import (
	"fmt"
	"github.com/eatMoreApple/regia"
	"os"
)

func main() {
	engine := regia.Default()
	engine.POST("/", func(ctx *regia.Context) {
		files, err := ctx.Request.Files()
		if err != nil {
			fmt.Println(err)
			return
		}
		file, err := files.Get("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		dst, _ := os.Create(file.Filename)
		defer dst.Close()
		file.Copy(dst)

		ctx.Response.String("upload success")
	})
	engine.Run(":8000")
}
```

可以通过`Context.SaveUploadFile`来对文件进行自定义操作。



#### Scan

```go
func (r *Request) Scan(scanner Scanner, v interface{}) error
```

实现`Scanner`, 获取当前的`*http.Request`对象进行操作



#### ScanJson

```go
func (r *Request) ScanJson(v interface{}) error
```

将request.Body以`JSON`格式解析到`v`上



#### ScanXml

```go
func (r *Request) ScanXml(v interface{}) error
```

将request.Body以`XML`格式解析到`v`上





### Response



#### SetStatus

```go
func (r *Response) SetStatus(code int) 
```

设置响应的状态码



#### SetHeader

```go
func (r *Response) SetHeader(key, value string)
```

设置响应头



#### SetCookie

设置响应`COOKIE`



#### Render

```go
func (r *Response) Render(render Render, data interface{}) error
```

渲染响应对象



#### Json

```go
func (r *Response) Json(data interface{}) error
```

将`data`转换成`JSON`写入`Response`



#### Xml

```go
func (r *Response) Xml(data interface{}) error
```

将`data`转换成`xml`写入`Response`



#### String

```go
func (r *Response) String(format string, a ...interface{}) (int, error)
```



#### Html

```go
func (r *Response) Html(name string, data interface{}) error
```

模板渲染





#### Redirect

```go
func (r *Response) Redirect(code int, url string)
```





#### ServeFile

```go
func (r *Response) ServeFile(path string)
```





#### ServeContent

```go
func (r *Response) ServeContent(name string, modTime time.Time, content io.ReadSeeker)
```

