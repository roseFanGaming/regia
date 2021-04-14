package regia

// To Regia

import (
	"net/http"
	"strings"
)

const (
	FilePathParam = "FilePathParam"
	wildFilepath  = "*" + FilePathParam
)

// Engine is a collection of core components of the whole service
type Engine struct {
	// The branch are used to store the handler.
	// All handlers are going to register in the router
	*Branch

	// Router is a module used to register handle and distribute request
	Router Router

	// Response html render
	// default use regia.TemplateRender
	// reset it to other html render engine
	HtmlRender HtmlRender

	// global json serializer
	// default use `encoding/json`
	// reset it to other module
	JsonSerializer Serializer

	// global xml serializer
	// default use `encoding/xml`
	// reset it to other module
	XmlSerializer Serializer

	// Context.SaveUploadFile will call this interface
	// default save file to your local desk
	// reset it to your onw idea
	FileStorage FileStorage

	// Context Abort use
	// default do nothing
	// reset it to implement your idea
	Abort Exit

	// NotFoundHandle replies to the request with an HTTP 404 not found error.
	NotFoundHandle func(ctx *Context)

	// All requests will be intercepted by Interceptors
	// whatever route matched or not
	Interceptors HandleFuncGroup

	// Starter will run when the service starts
	// and it only run once
	Starters []Starter

	// Warehouse is used to store information
	Warehouse Warehouse

	// Mat multipart form memory size
	// default 32M
	MultipartFormMaxMemory int64
}

// register all handles to router
func (e *Engine) registerHandle() {
	for method, nodes := range e.methodsTree {
		for _, node := range nodes {
			e.Router.Insert(method, node.path, node.group)
		}
	}
}

// Setter for Engine.NotFoundHandle
func (e *Engine) SetNotFoundHandle(handle HandleFunc) {
	e.NotFoundHandle = handle
}

// Serve static files
func (e *Engine) Static(url, dir string, group ...HandleFunc) {
	if strings.Contains(url, "*") {
		panic("`url` should not have wildcards")
	}
	server := http.FileServer(http.Dir(dir))
	handle := func(ctx *Context) {
		ctx.Raw.Request.URL.Path = ctx.Request.Params.Get(FilePathParam).String()
		server.ServeHTTP(ctx.Raw.Writer, ctx.Raw.Request)
	}
	group = append(group, handle)
	if !strings.HasSuffix(url, FilePathParam) {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		url += wildFilepath
	}
	e.Handle(http.MethodGet, url, group...)
}

// Add interceptor to Engine
func (e *Engine) AddInterceptors(interceptors ...HandleFunc) {
	e.Interceptors = append(e.Interceptors, interceptors...)
}

// Add starter to Engine
func (e *Engine) AddStarter(starters ...Starter) {
	e.Starters = append(e.Starters, starters...)
}

// Call all starters of this engine
func (e *Engine) runStarter() {
	for _, starter := range e.Starters {
		starter.Start(e)
	}
}

// Init engine
func (e *Engine) init() {
	e.registerHandle()
	e.runStarter()
}

// Start Listen and serve
func (e *Engine) Run(addr string) error {
	e.init()
	return http.ListenAndServe(addr, e)
}

// Handle input request
func (e *Engine) handleRequest(ctx *Context) {
	ctx.group = e.Interceptors
	if group, params := e.Router.Match(ctx.Raw.Request); group != nil {
		ctx.Request.Params = params
		ctx.group = append(ctx.group, group...)
	} else {
		ctx.group = append(ctx.group, e.NotFoundHandle)
	}
	ctx.start()
}

// ServeHTTP implement http.Handle
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := newContext(request, writer, e)
	e.handleRequest(ctx)
}

// Constructor for Engine
func New() *Engine {
	engine := &Engine{
		Router:                 make(HttpRouter),
		FileStorage:            &FileSystemStorage{},
		Branch:                 NewBranch(),
		JsonSerializer:         JsonSerializer{},
		XmlSerializer:          XmlSerializer{},
		HtmlRender:             TemplateRender{Template},
		Abort:                  exit{},
		NotFoundHandle:         HandleNotFound,
		Warehouse:              new(Data),
		MultipartFormMaxMemory: 32 << 20, // 32 MB
	}
	return engine
}

// Default Engine for use
func Default() *Engine {
	engine := New()
	engine.AddInterceptors(LogInterceptor)
	engine.AddStarter(&BannerStarter{Banner: Banner}, &UrlInfoStarter{})
	return engine
}

// Shortcut fot map[string]interface{}
type Map map[string]interface{}
