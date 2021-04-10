package regia

import (
	"errors"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	emptyValueError = errors.New("empty Value")
	emptyValue      = Value{err: emptyValueError}
)

type URLValue url.Values

func (u URLValue) Get(key string) Value {
	if u == nil {
		return emptyValue
	}
	vs := u[key]
	if vs == nil {
		return emptyValue
	}
	return Value{data: vs[0], err: nil}
}

func (u URLValue) GetAll(key string) Values {
	if u == nil {
		return Values{}
	}
	vs := u[key]
	var vas Values
	for _, i := range vs {
		vas = append(vas, newValue(i))
	}
	return vas
}

func (u URLValue) GetDefault(key, def string) Value {
	if u == nil {
		return newValue(def)
	}
	vs := u[key]
	if vs == nil {
		return newValue(def)
	}
	return newValue(vs[0])
}

type Value struct {
	data string
	err  error
}

func (v Value) Raw() (string, error) {
	return v.data, v.err
}

func (v Value) Int(def ...int) int {
	if v.IsValid() && !v.IsEmpty() {
		i, err := strconv.Atoi(v.data)
		if err != nil && def != nil {
			return def[len(def)-1]
		}
		return i
	}
	if def == nil {
		def = []int{0}
	}
	return def[0]
}

func (v Value) Int64(def ...int64) int64 {
	if v.IsValid() && !v.IsEmpty() {
		i, err := strconv.ParseInt(v.data, 10, 64)
		if err != nil && def != nil {
			return def[len(def)-1]
		}
		return i
	}
	if def == nil {
		def = []int64{0}
	}
	return def[0]
}

func (v Value) Float64(def ...float64) float64 {
	if v.IsValid() && !v.IsEmpty() {
		i, err := strconv.ParseFloat(v.data, 64)
		if err != nil && def != nil {
			return def[len(def)-1]
		}
		return i
	}
	if def == nil {
		def = []float64{0}
	}
	return def[0]
}

func (v Value) String(def ...string) string {
	if v.IsValid() && !v.IsEmpty() {
		return v.data
	}
	if def == nil {
		def = []string{""}
	}
	return def[0]
}

func (v Value) Time() (time.Time, error) {
	if v.IsValid() && !v.IsEmpty() {
		ts, err := strconv.ParseInt(v.data, 10, 64)
		if err == nil {
			return time.Unix(ts, 0), err
		}
	}
	return time.Time{}, v.err
}

func (v Value) ParseTime(layout string) (time.Time, error) {
	if v.IsValid() && !v.IsEmpty() {
		return time.Parse(layout, v.data)
	}
	return time.Time{}, v.err
}

func (v Value) IsEmpty() bool {
	return v.data == ""
}

func (v Value) IsValid() bool {
	return v.err == nil
}

func (v Value) Err() error {
	return v.err
}

type Values []Value

func newValue(data string) Value {
	return Value{data: data}
}

type Warehouse interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, exist bool)
}

type Data struct {
	item map[string]interface{}
	mu   sync.RWMutex
}

func (d *Data) Set(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.item == nil {
		d.item = make(map[string]interface{})
	}
	d.item[key] = value
}

func (d *Data) Get(key string) (value interface{}, exist bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, exist = d.item[key]
	return
}

func (d *Data) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.item = nil
}

func (d *Data) Reset() {
	d.Clear()
}
