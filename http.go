package regia

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	contentType         = "Content-Type"
	jsonContentType     = "application/json;charset=utf-8"
	textHtmlContentType = "text/html;charset=utf-8"
	textXmlContentType  = "text/xml;charset=utf-8"
)

const (
	MethodGet     = "Get"
	MethodPost    = "Post"
	MethodPut     = "Put"
	MethodPatch   = "Patch"
	MethodDelete  = "Delete"
	MethodHead    = "Head"
	MethodOptions = "Options"
	MethodTrace   = "Trace "
)

var httpMethods = [...]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodHead,
	http.MethodTrace,
}

// Http Request Method Mapping
var HttpRequestMethodMapping = map[string]string{
	MethodPost:    http.MethodPost,
	MethodGet:     http.MethodGet,
	MethodPut:     http.MethodPut,
	MethodPatch:   http.MethodPatch,
	MethodDelete:  http.MethodDelete,
	MethodHead:    http.MethodHead,
	MethodOptions: http.MethodOptions,
	MethodTrace:   http.MethodTrace,
}

// net/http raw value
type raw struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

type Request struct {
	*http.Request
	Context *Context
	Params  Params
}

func (r *Request) Query() URLValue {
	return URLValue(r.Request.URL.Query())
}

func (r *Request) Form() URLValue {
	return URLValue(r.Request.PostForm)
}

func (r *Request) Files() (Files, error) {
	req := r.Request
	if req.MultipartForm == multipartByReader {
		return nil, multipartReaderError
	}
	if req.MultipartForm == nil {
		if err := req.ParseMultipartForm(r.Context.Engine.MultipartFormMaxMemory); err != nil {
			return nil, err
		}
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		return req.MultipartForm.File, nil
	}
	return nil, http.ErrMissingFile
}

func (r *Request) Scan(scanner Scanner, v interface{}) error {
	return scanner.Scan(r.Context.Raw.Request, v)
}

func (r *Request) ScanJson(v interface{}) error {
	scanner := JsonScanner{Serializer: r.Context.Engine.JsonSerializer}
	return r.Scan(scanner, v)
}

func (r *Request) ScanXml(v interface{}) error {
	scanner := XmlScanner{Serializer: r.Context.Engine.JsonSerializer}
	return r.Scan(scanner, v)
}

func (r *Request) GetCookie(key string) (*http.Cookie, error) {
	return r.Request.Cookie(key)
}

type Response struct {
	Context *Context
	http.ResponseWriter
}

func (r *Response) SetStatus(code int) {
	r.ResponseWriter.WriteHeader(code)
}

func (r *Response) SetHeader(key, value string) {
	r.ResponseWriter.Header().Set(key, value)
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.ResponseWriter, cookie)
}

func (r *Response) Render(render Render, data interface{}) error {
	return render.Render(r.ResponseWriter, data)
}

func (r *Response) Json(data interface{}) error {
	render := JsonRender{Serializer: r.Context.Engine.JsonSerializer}
	return r.Render(render, data)
}

func (r *Response) String(format string, a ...interface{}) (int, error) {
	text := fmt.Sprintf(format, a...)
	writeContentType(r.Context.Raw.Writer, textHtmlContentType)
	return r.Context.Raw.Writer.Write([]byte(text))
}

func (r *Response) Xml(data interface{}) error {
	render := XmlRender{Serializer: r.Context.Engine.XmlSerializer}
	return r.Render(render, data)
}

func (r *Response) Html(name string, data interface{}) error {
	return r.Context.Engine.HtmlRender.Render(r.Context.Raw.Writer, name, data)
}

// Shortcut for http.Redirect
func (r *Response) Redirect(code int, url string) {
	http.Redirect(r.ResponseWriter, r.Context.Raw.Request, url, code)
}

// Shortcut for http.ServeFile
func (r *Response) ServeFile(path string) {
	http.ServeFile(r.ResponseWriter, r.Context.Raw.Request, path)
}

// Shortcut for http.ServeContent
func (r *Response) ServeContent(name string, modTime time.Time, content io.ReadSeeker) {
	http.ServeContent(r.ResponseWriter, r.Context.Raw.Request, name, modTime, content)
}

func writeContentType(writer http.ResponseWriter, cT string) {
	writer.Header().Del(contentType)
	writer.Header().Set(contentType, cT)
}
