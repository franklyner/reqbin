package reqbin

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	TagName = "param"
)

var TimeFormats = []string{time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.RFC822, time.RFC822Z, time.RFC850, time.UnixDate, "2006-01-02"}

/*
UnmarshallRequestForm Allows to parse an http.Request and map the corresponding values
in a very similar way as json.Unmarshall() does.

It will check if request.ParseForm() has been executed. If not, it executes it. It will then
use the request.FormValue(name string) to fetch the values. So, it will work for both:
query parameters and multi-part form values.

Usage:
Given the struct

	type TestStruct struct {
		Name    string    `param:"name"`
		IsCool  bool      `param:"is_cool"`
		Counter int       `param:"counter"`
		Start   time.Time `param:"start"`
	}

and the request query string:
name=Joe&is_cool=true&counter=1&start=20023-01-02T15:04:05Z

it will use the param tag to populate the struct fields accordingly. All basic types are supported as
well as time.Time (all formats defined in RFC constants in the time package, plus "2006-01-02")
*/
func UnmarshallRequestForm(request *http.Request, val any) error {
	t, err := validateAndGetType(val)
	if err != nil {
		return err
	}
	// ensure request.ParseForm has been done
	if request.Form == nil && request.PostForm == nil {
		if err := request.ParseForm(); err != nil {
			return fmt.Errorf("error parsing request form: %w", err)
		}
	}
	paramToFieldname := getListOfParamNames(t)
	s := reflect.ValueOf(val).Elem()
	for p, fn := range paramToFieldname {
		if err := setFieldValue(s, p, fn, request); err != nil {
			return fmt.Errorf("error setting value for param %s: %w", p, err)
		}
	}
	return nil
}

func setFieldValue(s reflect.Value, param string, fieldName string, request *http.Request) error {
	f := s.FieldByName(fieldName)
	if !f.CanSet() {
		return nil
	}
	requestValue := request.FormValue(param)
	if requestValue == "" {
		return nil
	}

	var err error
	requestValue, err = url.QueryUnescape(requestValue)
	if err != nil {
		return err
	}

	switch f.Kind() {
	case reflect.String:
		f.SetString(requestValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int
		if i, err = strconv.Atoi(requestValue); err == nil {
			f.SetInt(int64(i))
		}
	case reflect.Bool:
		var b bool
		if b, err = strconv.ParseBool(requestValue); err == nil {
			f.SetBool(b)
		}
	case reflect.Float32:
		var fl float64
		if fl, err = strconv.ParseFloat(requestValue, 32); err == nil {
			f.SetFloat(fl)
		}
	case reflect.Float64:
		var fl float64
		if fl, err = strconv.ParseFloat(requestValue, 64); err == nil {
			f.SetFloat(fl)
		}
	case reflect.Struct:
		switch f.Interface().(type) {
		case time.Time:
			var pt time.Time
			found := false
			for _, format := range TimeFormats {
				var eri error
				pt, eri = time.Parse(format, requestValue)
				if eri == nil {
					found = true
					break
				}
			}
			if found {
				f.Set(reflect.ValueOf(pt))
			} else {
				err = fmt.Errorf("invalid time format for parameter %s", fieldName)
			}
		}

	default:
		err = fmt.Errorf("unsupported type: %s", f.Kind())
	}
	return err
}

func validateAndGetType(val any) (reflect.Type, error) {
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind != reflect.Pointer {
		return v.Type(), errors.New("val is not a pointer to struct")
	}
	v = v.Elem()
	kind = v.Kind()
	if kind != reflect.Struct {
		return v.Type(), errors.New("val is not a struct to struct")
	}
	return v.Type(), nil
}

func getListOfParamNames(t reflect.Type) map[string]string {
	fieldCount := t.NumField()
	params := make(map[string]string)
	for i := 0; i < fieldCount; i++ {
		f := t.Field(i)
		t := f.Tag
		v := t.Get(TagName)
		v = strings.Split(v, ",")[0]
		params[v] = f.Name
	}
	return params
}
