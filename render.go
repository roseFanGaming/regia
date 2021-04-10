package regia

import (
	"html/template"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct {
	Serializer Serializer
}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	data, err := j.Serializer.Marshal(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

type XmlRender struct {
	Serializer Serializer
}

func (x XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, textXmlContentType)
	data, err := x.Serializer.Marshal(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

type HtmlRender interface {
	Render(writer http.ResponseWriter, name string, data interface{}) error
}

type TemplateRender struct {
	Template *template.Template
}

func (h TemplateRender) Render(writer http.ResponseWriter, name string, data interface{}) error {
	writeContentType(writer, textHtmlContentType)
	return h.Template.ExecuteTemplate(writer, name, data)
}

var Template = template.New("")
