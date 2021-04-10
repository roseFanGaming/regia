package regia

import (
	"encoding/json"
	"encoding/xml"
)

type Serializer interface {
	Unmarshal([]byte, interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

type JsonSerializer struct{}

func (j JsonSerializer) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

func (j JsonSerializer) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

type XmlSerializer struct{}

func (x XmlSerializer) Unmarshal(data []byte, v interface{}) error { return xml.Unmarshal(data, v) }

func (x XmlSerializer) Marshal(v interface{}) ([]byte, error) { return xml.Marshal(v) }
