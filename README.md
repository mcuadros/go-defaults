go-defaults [![Build Status](https://img.shields.io/github/workflow/status/mcuadros/go-defaults/Test.svg)](https://github.com/mcuadros/go-defaults/actions) [![GoDoc](http://godoc.org/github.com/mcuadros/go-defaults?status.png)](https://pkg.go.dev/github.com/mcuadros/go-defaults) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-defaults.svg)](https://github.com/mcuadros/go-defaults/releases)
==============================

Enabling stuctures with defaults values using [struct tags](http://golang.org/pkg/reflect/#StructTag).

Installation
------------

The recommended way to install go-defaults

```
go get github.com/mcuadros/go-defaults
```

Examples
--------

A basic example:

```go
import (
    "fmt"
    "github.com/mcuadros/go-defaults"
    "time"
)

type ExampleBasic struct {
    Foo bool   `default:"true"` //<-- StructTag with a default key
    Bar string `default:"33"`
    Qux int8
    Dur time.Duration `default:"1m"`
}

func NewExampleBasic() *ExampleBasic {
    example := new(ExampleBasic)
    defaults.SetDefaults(example) //<-- This set the defaults values

    return example
}

...

test := NewExampleBasic()
fmt.Println(test.Foo) //Prints: true
fmt.Println(test.Bar) //Prints: 33
fmt.Println(test.Qux) //Prints:
fmt.Println(test.Dur) //Prints: 1m0s
```

Caveats
-------

At the moment, the way the default filler checks whether it should fill a struct field or not is by comparing the current field value with the corresponding zero value of that type. This has a subtle implication: the zero value set explicitly by you will get overriden by default value during `SetDefaults()` call. So if you need to set the field to container zero value, you need to set it explicitly AFTER setting the defaults.

Take the basic example in the above section and change it slightly:
```go

example := ExampleBasic{
    Bar: 0,
}
defaults.SetDefaults(example)
fmt.Println(example.Bar) //Prints: 33 instead of 0 (which is zero value for int)

example.Bar = 0 // set needed zero value AFTER applying defaults
fmt.Println(example.Bar) //Prints: 0

```

Pointer Set
-------

Pointer field struct is a tricky usage to avoid covering existed values. 

Take the basic example in the above section and change it slightly:
```go

type ExamplePointer struct {
    Foo *bool   `default:"true"` //<-- StructTag with a default key
    Bar *string `default:"example"`
    Qux *int    `default:"22"`
    Oed *int64  `default:"64"`
}

...

boolZero := false
stringZero := ""
intZero := 0
example := &ExamplePointer{
    Foo: &boolZero,
    Bar: &stringZero,
    Qux: &intZero,
}
defaults.SetDefaults(example)

fmt.Println(*example.Foo) //Prints: false (zero value `false` for bool but not for bool ptr)
fmt.Println(*example.Bar) //Prints: "" (print "" which set in advance, not "example" for default)
fmt.Println(*example.Qux) //Prints: 0 (0 instead of 22)
fmt.Println(*example.Oed) //Prints: 64 (64, because the ptr addr is nil when SetDefaults)

```

It's also a very useful feature for web application which default values are needed while binding request json.

For example:
```go
type ExamplePostBody struct {
    Foo *bool   `json:"foo" default:"true"` //<-- StructTag with a default key
    Bar *string `json:"bar" default:"example"`
    Qux *int    `json:"qux" default:"22"`
    Oed *int64  `json:"oed" default:"64"`
}
```

HTTP request seems like this:
```bash
curl --location --request POST ... \
... \
--header 'Content-Type: application/json' \
--data-raw '{
    "foo": false,
    "bar": "",
    "qux": 0
}'
```

Request handler:
```go
func PostExampleHandler(c *gin.Context) {
    var reqBody ExamplePostBody
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, nil)
        return
    }
    defaults.SetDefaults(&reqBody)

    fmt.Println(*reqBody.Foo) //Prints: false (zero value `false` for bool but not for bool ptr)
    fmt.Println(*reqBody.Bar) //Prints: "" (print "" which set in advance, not "example" for default)
    fmt.Println(*reqBody.Qux) //Prints: 0 (0 instead of 22, did not confused from whether zero value is in json or not)
    fmt.Println(*reqBody.Oed) //Prints: 64 (In this case "oed" is not in req json, so set default 64)

    ...
}
```

License
-------

MIT, see [LICENSE](LICENSE)
