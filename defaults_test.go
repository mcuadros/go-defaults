package defaults

import (
	"testing"
	"time"

	"bou.ke/monkey"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.

func Test(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", "2020-06-10 12:00:00")
		return t
	})

	TestingT(t)
}

type DefaultsSuite struct{}

var _ = Suite(&DefaultsSuite{})

type Parent struct {
	Children []Child
}

type Child struct {
	Name string
	Age  int `default:"10"`
}

type ChildPtr struct {
	Name *string
	Age  *int `default:"10"`
}

type ExampleBasic struct {
	Bool       bool    `default:"true"`
	Integer    int     `default:"33"`
	Integer8   int8    `default:"8"`
	Integer16  int16   `default:"16"`
	Integer32  int32   `default:"32"`
	Integer64  int64   `default:"64"`
	UInteger   uint    `default:"11"`
	UInteger8  uint8   `default:"18"`
	UInteger16 uint16  `default:"116"`
	UInteger32 uint32  `default:"132"`
	UInteger64 uint64  `default:"164"`
	String     string  `default:"foo"`
	Bytes      []byte  `default:"bar"`
	Float32    float32 `default:"3.2"`
	Float64    float64 `default:"6.4"`
	Struct     struct {
		Bool    bool `default:"true"`
		Integer int  `default:"33"`
	}
	Duration         time.Duration `default:"1s"`
	Children         []Child
	Second           time.Duration `default:"1s"`
	StringSlice      []string      `default:"[1,2,3,4]"`
	IntSlice         []int         `default:"[1,2,3,4]"`
	IntSliceSlice    [][]int       `default:"[[1],[2],[3],[4]]"`
	StringSliceSlice [][]string    `default:"[[1],[]]"`

	DateTime string `default:"{{date:1,-10,0}} {{time:1,-5,10}}"`

	BoolPtr     *bool          `default:"false"`
	IntPtr      *int           `default:"33"`
	Int8Ptr     *int8          `default:"8"`
	Int16Ptr    *int16         `default:"16"`
	Int32Ptr    *int32         `default:"32"`
	Int64Ptr    *int64         `default:"64"`
	UIntPtr     *uint          `default:"11"`
	UInt8Ptr    *uint8         `default:"18"`
	UInt16Ptr   *uint16        `default:"116"`
	UInt32Ptr   *uint32        `default:"132"`
	UInt64Ptr   *uint64        `default:"164"`
	Float32Ptr  *float32       `default:"3.2"`
	Float64Ptr  *float64       `default:"6.4"`
	DurationPtr *time.Duration `default:"1s"`
	SecondPtr   *time.Duration `default:"1s"`
	StructPtr   *struct {
		Bool    bool `default:"true"`
		Integer *int `default:"33"`
	}
	PtrStructPtr **struct {
		Bool    bool `default:"false"`
		Integer *int `default:"33"`
	}
	ChildrenPtr         []*ChildPtr
	PtrChildrenPtr      *[]*ChildPtr
	PtrPtrChildrenPtr   **[]*ChildPtr
	PtrStringSliceNoTag *[]string
	PtrStringSlice      *[]string   `default:"[1,2,3,4]"`
	PtrIntSlice         *[]int      `default:"[1,2,3,4]"`
	PtrIntSliceSlice    *[][]int    `default:"[[1],[2],[3],[4]]"`
	PtrStringSliceSlice *[][]string `default:"[[1],[]]"`
	Float64PtrNoTag     *float64
}

func (s *DefaultsSuite) TestSetDefaultsBasic(c *C) {
	foo := &ExampleBasic{}
	SetDefaults(foo)

	s.assertTypes(c, foo)
}

type ExampleNested struct {
	Struct ExampleBasic
}

func (s *DefaultsSuite) TestSetDefaultsNested(c *C) {
	foo := &ExampleNested{}
	SetDefaults(foo)

	s.assertTypes(c, &foo.Struct)
}

func (s *DefaultsSuite) assertTypes(c *C, foo *ExampleBasic) {
	c.Assert(foo.Bool, Equals, true)
	c.Assert(foo.Integer, Equals, 33)
	c.Assert(foo.Integer8, Equals, int8(8))
	c.Assert(foo.Integer16, Equals, int16(16))
	c.Assert(foo.Integer32, Equals, int32(32))
	c.Assert(foo.Integer64, Equals, int64(64))
	c.Assert(foo.UInteger, Equals, uint(11))
	c.Assert(foo.UInteger8, Equals, uint8(18))
	c.Assert(foo.UInteger16, Equals, uint16(116))
	c.Assert(foo.UInteger32, Equals, uint32(132))
	c.Assert(foo.UInteger64, Equals, uint64(164))
	c.Assert(foo.String, Equals, "foo")
	c.Assert(string(foo.Bytes), Equals, "bar")
	c.Assert(foo.Float32, Equals, float32(3.2))
	c.Assert(foo.Float64, Equals, 6.4)
	c.Assert(foo.Struct.Bool, Equals, true)
	c.Assert(foo.Duration, Equals, time.Second)
	c.Assert(foo.Children, IsNil)
	c.Assert(foo.Second, Equals, time.Second)
	c.Assert(foo.StringSlice, DeepEquals, []string{"1", "2", "3", "4"})
	c.Assert(foo.IntSlice, DeepEquals, []int{1, 2, 3, 4})
	c.Assert(foo.IntSliceSlice, DeepEquals, [][]int{[]int{1}, []int{2}, []int{3}, []int{4}})
	c.Assert(foo.StringSliceSlice, DeepEquals, [][]string{[]string{"1"}, []string{}})
	c.Assert(foo.DateTime, Equals, "2020-08-10 12:55:10")
	c.Assert(*foo.BoolPtr, Equals, false)
	c.Assert(*foo.IntPtr, Equals, 33)
	c.Assert(*foo.Int8Ptr, Equals, int8(8))
	c.Assert(*foo.Int16Ptr, Equals, int16(16))
	c.Assert(*foo.Int32Ptr, Equals, int32(32))
	c.Assert(*foo.Int64Ptr, Equals, int64(64))
	c.Assert(*foo.UIntPtr, Equals, uint(11))
	c.Assert(*foo.UInt8Ptr, Equals, uint8(18))
	c.Assert(*foo.UInt16Ptr, Equals, uint16(116))
	c.Assert(*foo.UInt32Ptr, Equals, uint32(132))
	c.Assert(*foo.UInt64Ptr, Equals, uint64(164))
	c.Assert(*foo.Float32Ptr, Equals, float32(3.2))
	c.Assert(*foo.Float64Ptr, Equals, 6.4)
	c.Assert(*foo.DurationPtr, Equals, time.Second)
	c.Assert(*foo.SecondPtr, Equals, time.Second)
	c.Assert(foo.StructPtr.Bool, Equals, true)
	c.Assert(*foo.StructPtr.Integer, Equals, 33)
	c.Assert((*foo.PtrStructPtr).Bool, Equals, false)
	c.Assert(*(*foo.PtrStructPtr).Integer, Equals, 33)
	c.Assert(foo.ChildrenPtr, IsNil)
	c.Assert(*foo.PtrChildrenPtr, IsNil)
	c.Assert(**foo.PtrPtrChildrenPtr, IsNil)
	c.Assert(*foo.PtrStringSliceNoTag, IsNil)
	c.Assert(*foo.PtrStringSlice, DeepEquals, []string{"1", "2", "3", "4"})
	c.Assert(*foo.PtrIntSlice, DeepEquals, []int{1, 2, 3, 4})
	c.Assert(*foo.PtrIntSliceSlice, DeepEquals, [][]int{[]int{1}, []int{2}, []int{3}, []int{4}})
	c.Assert(*foo.PtrStringSliceSlice, DeepEquals, [][]string{[]string{"1"}, []string{}})
	c.Assert(foo.Float64PtrNoTag, IsNil)
}

func (s *DefaultsSuite) TestSetDefaultsWithValues(c *C) {
	foo := &ExampleBasic{
		Integer:  55,
		UInteger: 22,
		Float32:  9.9,
		String:   "bar",
		Bytes:    []byte("foo"),
		Children: []Child{{Name: "alice"}, {Name: "bob", Age: 2}},
	}

	intzero := 0
	foo.IntPtr = &intzero

	ageZero := 0
	childPtr := &ChildPtr{Age: &ageZero}
	foo.ChildrenPtr = append(foo.ChildrenPtr, childPtr)

	SetDefaults(foo)

	c.Assert(foo.Integer, Equals, 55)
	c.Assert(foo.UInteger, Equals, uint(22))
	c.Assert(foo.Float32, Equals, float32(9.9))
	c.Assert(foo.String, Equals, "bar")
	c.Assert(string(foo.Bytes), Equals, "foo")
	c.Assert(foo.Children[0].Age, Equals, 10)
	c.Assert(foo.Children[1].Age, Equals, 2)
	c.Assert(*foo.ChildrenPtr[0].Age, Equals, 0)
	c.Assert(foo.ChildrenPtr[0].Name, IsNil)

	c.Assert(*foo.IntPtr, Equals, 0)
}

func (s *DefaultsSuite) BenchmarkLogic(c *C) {
	for i := 0; i < c.N; i++ {
		foo := &ExampleBasic{}
		SetDefaults(foo)
	}
}
