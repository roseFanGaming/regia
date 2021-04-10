package regia

import (
	"net/http"
)

type Exit interface{ Exit(ctx *Context) }

type exit struct{}

// Do nothing
func (e exit) Exit(*Context) {}

type Context struct {
	Raw      *raw
	Data     *Data
	group    HandleFuncGroup
	index    int
	Engine   *Engine
	Request  *Request
	Response *Response
}

func (c *Context) init() {
	c.Data = new(Data)
	_ = c.Raw.Request.ParseForm()
}

func (c *Context) start() {
	defer c.recover()
	c.Next()
}

func (c *Context) Next() {
	c.index++
	for c.index <= len(c.group) {
		handle := c.group[c.index-1]
		handle(c)
		c.index++
	}
}

func (c *Context) recover() {
	if rec := recover(); rec != nil {
		if e, ok := rec.(Exit); !ok {
			panic(rec)
		} else {
			e.Exit(c)
		}
	}
}

func (c *Context) Abort() { c.AbortWith(c.Engine.Abort) }

func (c *Context) AbortWith(exit Exit) { panic(exit) }

// Make http.ResponseWriter as http.Flusher
func (c *Context) Flusher() http.Flusher { return c.Raw.Writer.(http.Flusher) }

// implement your own idea with it
func (c *Context) SaveUploadFile(filer *File, path string) error {
	return c.Engine.FileStorage.Save(filer, path)
}

func (c *Context) setWithRaw(req *http.Request, writer http.ResponseWriter, engine *Engine) {
	c.Raw = &raw{Request: req, Writer: writer}
	c.Engine = engine
	c.Request = &Request{Context: c, Request: req}
	c.Response = &Response{Context: c, ResponseWriter: writer}
}

func newContext(req *http.Request, writer http.ResponseWriter, engine *Engine) *Context {
	ctx := new(Context)
	ctx.setWithRaw(req, writer, engine)
	ctx.init()
	return ctx
}
