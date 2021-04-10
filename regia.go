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

type Engine struct {
	*Branch
	Router                 Router
	HtmlRender             HtmlRender
	JsonSerializer         Serializer
	XmlSerializer          Serializer
	FileStorage            FileStorage
	Abort                  Exit
	NotFoundHandle         func(ctx *Context)
	Interceptors           HandleFuncGroup
	Starters               []Starter
	Warehouse              Warehouse
	MultipartFormMaxMemory int64
}

func (e *Engine) registerHandle() {
	for method, nodes := range e.methodsTree {
		for _, node := range nodes {
			e.Router.Insert(method, node.path, node.group)
		}
	}
}

func (e *Engine) SetNotFoundHandle(handle HandleFunc) {
	e.NotFoundHandle = handle
}

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

func (e *Engine) AddInterceptors(interceptors ...HandleFunc) {
	e.Interceptors = append(e.Interceptors, interceptors...)
}

func (e *Engine) AddStarter(starters ...Starter) {
	e.Starters = append(e.Starters, starters...)
}

func (e *Engine) runStarter() {
	for _, starter := range e.Starters {
		starter.Start(e)
	}
}

func (e *Engine) init() {
	e.registerHandle()
	e.runStarter()
}

func (e *Engine) Run(addr string) error {
	e.init()
	return http.ListenAndServe(addr, e)
}

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

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := newContext(request, writer, e)
	e.handleRequest(ctx)
}

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

func Default() *Engine {
	engine := New()
	engine.AddInterceptors(LogInterceptor)
	engine.AddStarter(&BannerStarter{Banner: Banner}, &UrlInfoStarter{})
	return engine
}

type Map map[string]interface{}
