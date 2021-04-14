# regia

Regia is a web framework written with golang ! 

It is simple, helpful and with high performance. Build your own idea with it !

## Installation

Golang version 1.11 + required

```shell
go get github.com/eatMoreApple/regia
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




