# regia

Regia is a web framework written with golang ! 

It is simple, helpful and easy to use. Build your own idea with it !

## Installation

Golang version 1.11 + required

```shell
go get https://github.com/eatMoreApple/regia
```



## Quick Start

```sh
$ touch main.go
# add all following codes into main.go
```

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
	engine.GET("/", func(ctx *regia.Context) {
		ctx.Response.Json(regia.Map{"hello": "world"})
	})
	engine.Run(":8000")
}
```

```shell
$ go run main.go
# open your brower and visit `localhost:8000/`
```



## Apis



### Routers



#### Register Router

```go
package main

import (
	"github.com/eatMoreApple/regia"
	"net/http"
)

func hello(ctx *regia.Context) {
	ctx.Response.Json(regia.Map{"hello": "world"})
}

// middleware 
// check request is allowed
func middleware(ctx *regia.Context) {
	if token := ctx.Request.Query().Get("token"); !token.IsValid() {
		ctx.Response.SetStatus(http.StatusForbidden)
		ctx.Response.String("permission denied")
		ctx.Abort() // exit
	}
  ctx.Next() // continue
}

func main() {
	engine := regia.Default()
	engine.GET("/", middleware, hello)
	engine.POST("/", hello)
	engine.PUT("/", hello)
	engine.PATCH("/", hello)
	engine.DELETE("/", hello)
	engine.OPTIONS("/", hello)
	//engine.Any("/", hello)
	engine.Run(":8000")
}
```



#### Register Struct Router

 ```go
package main

import (
	"github.com/eatMoreApple/regia"
)

type Dispatcher struct{}

func (d Dispatcher) Get(ctx *regia.Context) {
	ctx.Response.String("%s", ctx.Request.Method())
}

func (d Dispatcher) Post(ctx *regia.Context) {
	ctx.Response.String("%s", ctx.Request.Method())
}

func (d Dispatcher) Put() {} // this handle won't be register in

func main() {
	engine := regia.Default()
	engine.BindMethod("/", &Dispatcher{})
	engine.Run(":8000")
}
 ```

```shell
[REGIA URL INFO]   POST        1 handlers   /
[REGIA URL INFO]   GET         1 handlers   /
```



#### Register Struct With Custom Rule

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

type NameDispatcher struct{}

func (*NameDispatcher) Add(ctx *regia.Context) {
	ctx.Response.String("add page")
}

func main() {
	engine := regia.Default()
	engine.Bind("/", &NameDispatcher{}, map[string]string{"Add": "get"})
	engine.Run(":8000")
}
```



#### Dynamic Url

```go
package main

import "github.com/eatMoreApple/regia"

func main() {
	engine := regia.Default()
  
	engine.GET("/:id", func(ctx *regia.Context) {
		if id := ctx.Request.Params.Get("id"); id.IsValid() {
			ctx.Response.String(id.String())
		}
	})
  
	engine.Run(":8000")
}
```



### Request

#### Query Params

```go
package main

import (
	"fmt"
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
  
	engine.GET("/", func(ctx *regia.Context) {
		name := ctx.Request.Query().Get("name")
		age := ctx.Request.Query().GetDefault("age", "1")
		bobbies := ctx.Request.Query().GetAll("hobby")
		fmt.Println(name.String())
		fmt.Println(age.String())
		fmt.Println(bobbies)
		ctx.Response.String("ok")
		//get http://localhost:8000/?name=ivy&hobby=dance&hobby=basketball&hobby=sing
	})
  
	engine.Run(":8000")
}
```

```shell
ivy
1
[{dance <nil>} {basketball <nil>} {sing <nil>}]
```



#### Form 

```go
package main

import (
	"fmt"
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
  
	engine.POST("/", func(ctx *regia.Context) {
		name := ctx.Request.Form().Get("name")
		age := ctx.Request.Form().GetDefault("age", "1")
		bobbies := ctx.Request.Form().GetAll("hobby")
		ctx.Response.String("ok")
		fmt.Println(name)
		fmt.Println(age)
		fmt.Println(bobbies)
		// post http://localhost:8000/?name=ivy&hobby=dance&hobby=basketball&hobby=sing
	})
  
	engine.Run(":8000")
}
```

```shell
{ivy <nil>}
{1 <nil>}
[{dance <nil>} {basketball <nil>} {sing <nil>}]
```



### Return Response

**Json**	

```go
ctx.Response.Json(regia.Map{"hello": "world"})
```

Reset `JsonSerializer` with `Engine `to change default  behavior 

Default `JsonSerializer ` `encoding/json`. 

Use `jsoniter` to replace 

```go
package main

import (
	"github.com/eatMoreApple/regia"
	"github.com/json-iterator/go"
)

func main() {
	engine := regia.Default()
	engine.JsonSerializer = jsoniter.ConfigCompatibleWithStandardLibrary
}
```



**Xml**	

```go
type item struct { Name string `xml:"name"`}

ctx.Response.Xml(&item{Name: "ivy"})
```

Reset `XmlSerializer ` with `Engine ` to change default  behavior 

Default `XmlSerializer ` `encoding/xml`



**Html**

```go
package main

import (
	"github.com/eatMoreApple/regia"
	"html/template"
)

func main() {
	regia.Template = template.Must(template.ParseFiles("1.html"))
	engine := regia.Default()
	engine.GET("/", func(ctx *regia.Context) {
		ctx.Response.Html("1.html", nil)
	})
	engine.Run(":8000")
}
```

Reset `HtmlRender ` with `Engine ` to change default  behavior 

Default `HtmlRender ` `regia.TemplateRender`



#### Branch

```go
package main

import (
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()
	branch := regia.NewBranch()
	branch.GET("/login", func(ctx *regia.Context) {
		ctx.Response.String("login success")
	})
	engine.Include("/user", branch)
	engine.Run(":8000")
}
```

```shell
[REGIA URL INFO]   GET         1 handlers   /user/login
```



#### Static Files

```go
engine.Static("/static/", "static") // can add middleware here
```

```shell
[REGIA URL INFO]   GET         1 handlers   /static/*FilePathParam
```



#### Midlleware

```go
package main

import (
	"fmt"
	"github.com/eatMoreApple/regia"
)

func main() {
	engine := regia.Default()

	// global middleware
	engine.Use(func(ctx *regia.Context) { fmt.Println("middleware1") })

	branch := regia.NewBranch()

	// branch middleware
	branch.Use(func(ctx *regia.Context) { fmt.Println("middleware3") })

	// your handle
	login := func(ctx *regia.Context) {
		fmt.Println("end")
		ctx.Response.String("login page")
	}

	// handle middleware
	branch.GET("/login", func(ctx *regia.Context) {
		fmt.Println("middleware4")
	}, login)

	// branch middleware
	// this middleware will be called at first of branch middleware
	engine.Include("/branch", branch)

	engine.Run(":8000")
}
```

```shell
middleware1
middleware2
middleware3
middleware4
end
```





