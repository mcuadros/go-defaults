go-defaults [![Build Status](https://travis-ci.org/mcuadros/go-defaults.png?branch=master)](https://travis-ci.org/mcuadros/go-defaults) [![GoDoc](http://godoc.org/github.com/mcuadros/go-defaults?status.png)](http://godoc.org/github.com/mcuadros/go-defaults) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-defaults.svg)](https://github.com/mcuadros/go-defaults/releases)
==============================

Enabling stuctures with defaults values using [struct tags](http://golang.org/pkg/reflect/#StructTag).

Installation
------------

The recommended way to install go-defaults

```
go get gopkg.in/mcuadros/go-defaults.v1
```

Examples
--------

A basic example:

```go
import (
    "fmt"
    "github.com/mcuadros/go-defaults"
)

type ExampleBasic struct {
    Foo bool   `default:"true"` //<-- StructTag with a default key
    Bar string `default:"33"`
    Qux int8
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

```

License
-------

MIT, see [LICENSE](LICENSE)
