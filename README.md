go-defaults [![Build Status](https://travis-ci.org/mcuadros/go-defaults.png?branch=master)](https://travis-ci.org/mcuadros/go-defaults) [![GoDoc](https://godoc.org/github.com/mcuadros/go-defaults?status.png)](http://godoc.org/github.com/mcuadros/go-defaults)
==============================

This library allow to define a default value to any struct, this is made thanks to [struct tags](http://golang.org/pkg/reflect/#StructTag).

> A StructTag is the tag string in a struct field.

> By convention, tag strings are a concatenation of optionally space-separated key:"value" pairs. Each key is a non-empty string consisting of non-control characters other than space (U+0020 ' '), quote (U+0022 '"'), and colon (U+003A ':'). Each value is quoted using U+0022 '"' characters and Go string literal syntax.


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
import "github.com/mcuadros/go-defaults"
import "fmt"

type ExampleBasic struct {
    Foo bool   `default:"true"` //<-- StructTag with a default key
    Bar string `default:"33"`
    Qux int8
}

func NewExampleBasic() *ExampleBasic {
    example := new(ExampleBasic)
    SetDefaults(example) //<-- This set the defaults values

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
