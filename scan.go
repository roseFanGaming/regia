package regia

import (
	"bytes"
	"net/http"
)

type Scanner interface {
	Scan(req *http.Request, data interface{}) error
}

type JsonScanner struct {
	Serializer Serializer
}

func (j JsonScanner) Scan(req *http.Request, v interface{}) error {
	buffer := &bytes.Buffer{}
	if _, err := buffer.ReadFrom(req.Body); err != nil {
		return err
	}
	return j.Serializer.Unmarshal(buffer.Bytes(), v)
}

type XmlScanner struct {
	Serializer Serializer
}

func (x XmlScanner) Scan(req *http.Request, v interface{}) error {
	buffer := &bytes.Buffer{}
	if _, err := buffer.ReadFrom(req.Body); err != nil {
		return err
	}
	return x.Serializer.Unmarshal(buffer.Bytes(), v)
}
