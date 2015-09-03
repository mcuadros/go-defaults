package defaults

import (
	"fmt"
	"reflect"

	. "gopkg.in/check.v1"
)

type FillerSuite struct{}

var _ = Suite(&FillerSuite{})

type FixtureTypeInt int

func (s *FillerSuite) TestFuncByNameIsEmpty(c *C) {
	calledA := false
	calledB := false

	f := &Filler{
		FuncByName: map[string]fillerFunc{
			"Foo": func(field *fieldData) {
				calledA = true
			},
		},
		FuncByKind: map[reflect.Kind]fillerFunc{
			reflect.Int: func(field *fieldData) {
				calledB = true
			},
		},
	}

	f.Fill(&struct{ Foo int }{})
	c.Assert(calledA, Equals, true)
	c.Assert(calledB, Equals, false)
}

func (s *FillerSuite) TestFuncByTypeIsEmpty(c *C) {
	calledA := false
	calledB := false

	t := GetTypeHash(reflect.TypeOf(new(FixtureTypeInt)))
	f := &Filler{
		FuncByType: map[TypeHash]fillerFunc{
			t: func(field *fieldData) {
				calledA = true
			},
		},
		FuncByKind: map[reflect.Kind]fillerFunc{
			reflect.Int: func(field *fieldData) {
				calledB = true
			},
		},
	}

	f.Fill(&struct{ Foo FixtureTypeInt }{})
	c.Assert(calledA, Equals, true)
	c.Assert(calledB, Equals, false)
}

func (s *FillerSuite) TestFuncByKindIsNotEmpty(c *C) {
	called := false
	f := &Filler{FuncByKind: map[reflect.Kind]fillerFunc{
		reflect.Int: func(field *fieldData) {
			called = true
		},
	}}

	f.Fill(&struct{ Foo int }{Foo: 42})
	c.Assert(called, Equals, false)
}

func (s *FillerSuite) TestFuncByKindSlice(c *C) {
	fmt.Println(GetTypeHash(reflect.TypeOf(new([]string))))
}

func (s *FillerSuite) TestFuncByKindTag(c *C) {
	var called string
	f := &Filler{Tag: "foo", FuncByKind: map[reflect.Kind]fillerFunc{
		reflect.Int: func(field *fieldData) {
			called = field.TagValue
		},
	}}

	f.Fill(&struct {
		Foo int `foo:"qux"`
	}{})
	c.Assert(called, Equals, "qux")
}

func (s *FillerSuite) TestFuncByKindIsEmpty(c *C) {
	called := false
	f := &Filler{FuncByKind: map[reflect.Kind]fillerFunc{
		reflect.Int: func(field *fieldData) {
			called = true
		},
	}}

	f.Fill(&struct{ Foo int }{})
	c.Assert(called, Equals, true)
}
