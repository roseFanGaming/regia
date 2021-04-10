package regia

import (
	"fmt"
	"net/http"
	"time"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

var logTitle = formatColor("[REGIA LOG]", colorGreen)

type HandleFunc func(ctx *Context)

type HandleFuncGroup []HandleFunc

func HandleWithValue(key string, value interface{}) HandleFunc {
	return func(ctx *Context) { ctx.Data.Set(key, value) }
}

func HandleNotFound(ctx *Context) { http.NotFound(ctx.Raw.Writer, ctx.Raw.Request) }

func LogInterceptor(ctx *Context) {
	start := time.Now()
	ctx.Next()
	endTime := time.Since(start)
	startTimeStr := formatColor(start.Format(defaultTimeFormat), colorYellow)
	method := formatColor(fmt.Sprintf("[METHOD:%s]", ctx.Raw.Request.Method), colorBlue)
	path := formatColor(fmt.Sprintf("[PATH:%s]", ctx.Raw.Request.URL.Path), 96) // #02F3F3
	addr := formatColor(fmt.Sprintf("[Addr:%s]", ctx.Raw.Request.RemoteAddr), 97)
	end := formatColor(endTime.String(), colorMagenta)
	// 2006-01-02 15:04:05     [METHOD:GET]     [Addr:127.0.0.1:49453]      [PATH:/name]
	fmt.Printf("%-20s %-32s %-20s %-28s %-35s  %-20s\n", logTitle, startTimeStr, end, method, addr, path)
}
