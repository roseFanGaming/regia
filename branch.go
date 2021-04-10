package regia

import (
	"net/http"
	"reflect"
)

type handleNode struct {
	path  string
	group HandleFuncGroup
}

type Branch struct {
	methodsTree map[string][]*handleNode
	middleware  HandleFuncGroup
	prefix      string
}

func (b *Branch) Use(group ...HandleFunc) { b.middleware = append(b.middleware, group...) }

func (b *Branch) SetPrefix(path string) { b.prefix = path }

func (b *Branch) GET(path string, group ...HandleFunc) {
	b.Handle(http.MethodGet, path, group...)
}

func (b *Branch) POST(path string, group ...HandleFunc) {
	b.Handle(http.MethodPost, path, group...)
}

func (b *Branch) PUT(path string, group ...HandleFunc) {
	b.Handle(http.MethodPut, path, group...)
}

func (b *Branch) PATCH(path string, group ...HandleFunc) {
	b.Handle(http.MethodPatch, path, group...)
}

func (b *Branch) DELETE(path string, group ...HandleFunc) {
	b.Handle(http.MethodDelete, path, group...)
}

func (b *Branch) HEAD(path string, group ...HandleFunc) {
	b.Handle(http.MethodHead, path, group...)
}

func (b *Branch) OPTIONS(path string, group ...HandleFunc) {
	b.Handle(http.MethodOptions, path, group...)
}

func (b *Branch) Any(path string, group ...HandleFunc) {
	for _, method := range httpMethods {
		b.Handle(method, path, group...)
	}
}

func (b *Branch) Handle(method, path string, group ...HandleFunc) {
	group = append(b.middleware, group...)
	path = b.prefix + path
	n := &handleNode{path: path, group: group}
	b.methodsTree[method] = append(b.methodsTree[method], n)
}

func (b *Branch) Include(prefix string, branch *Branch) {
	for method, nodes := range branch.methodsTree {
		for _, node := range nodes {
			b.Handle(method, prefix+node.path, node.group...)
		}
	}
}

func (b *Branch) Bind(path string, v interface{}, mappings ...map[string]string) {
	for _, mapping := range mappings {
		cleanedMapping := getCleanedRequestMapping(mapping)
		value := reflect.ValueOf(v)
		for handleName, methodName := range cleanedMapping {
			if method := value.MethodByName(handleName); method.IsValid() {
				if handle, ok := method.Interface().(func(ctx *Context)); ok {
					b.Handle(methodName, path, handle)
				}
			}
		}
	}
}

func (b *Branch) BindMethod(path string, v interface{}, mappings ...map[string]string) {
	mappings = append(mappings, HttpRequestMethodMapping)
	b.Bind(path, v, mappings...)
}

func NewBranch() *Branch {
	return &Branch{methodsTree: make(map[string][]*handleNode)}
}
