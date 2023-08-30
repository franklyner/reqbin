# reqbin
## Welcome to reqbin

reqbin is short for request binding. It aims at solving the problem to parse a URL query or a multi-part form POST into a struct.
The goal is to provide a similar experience than unmarshalling a JSON byte array into a struct.

Usage:
```
go get "github.com/franklyner/reqbin"

import (
  "github.com/franklyner/reqbin"
)

type MyStruct {
		Name    string    `param:"name"`
		IsCool  bool      `param:"is_cool"`
		Counter int       `param:"counter"`
		Start   time.Time `param:"start"`
}
...
req := ... // the http.Request being processed
s := &MyStruct{}
if err := reqbin.UnmarshallRequestForm(req, s); err != nil {
  //TODO
}
```
If the query of the request contained the parameters with the names defined in the struct tags, then fields of s will be populated with the corresponding values.

