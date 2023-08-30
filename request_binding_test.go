package main

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type TestStruct struct {
	Name    string    `param:"name"`
	IsCool  bool      `param:"is_cool"`
	Counter int       `param:"counter"`
	Start   time.Time `param:"start"`
}

func TestValidateAndGetType(t *testing.T) {
	s := struct{ n int }{3}
	errVals := []any{0, "test", s, make(map[string]string)}
	for _, v := range errVals {
		_, err := validateAndGetType(v)
		if err == nil {
			t.Errorf("value was supposed to fail: %+v", v)
		}
	}
	_, err := validateAndGetType(&s)
	if err != nil {
		t.Errorf("was supposed to pass: %s", err.Error())
	}
}

func TestGetParamsAndFields(t *testing.T) {
	s := TestStruct{}
	m := getListOfParamNames(reflect.ValueOf(s).Type())
	if len(m) != 3 {
		t.Errorf("didn't get expected map: %+v", m)
	}
}

func TestFull(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "http://something.com?name=Joe&is_cool=true&counter=1&start=20023-01-02T15:04:05Z", bytes.NewReader([]byte{}))
	s := &TestStruct{}
	err := UnmarshallRequestForm(r, s)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("s now looks like this: %+v", s)
}
