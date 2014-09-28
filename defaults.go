package defaults

import (
	"reflect"
	"strconv"
)

// Applies the default values to the struct object, the struct type must have
// the StructTag with name "default" and the directed value.
//
// Usage
//     type ExampleBasic struct {
//         Foo bool   `default:"true"`
//         Bar string `default:"33"`
//         Qux int8
//     }
//
//      foo := &ExampleBasic{}
//      SetDefaults(foo)
func SetDefaults(variable interface{}) {
	getDefaultFiller().Fill(variable)
}

var defaultFiller *Filler = nil

func getDefaultFiller() *Filler {
	if defaultFiller == nil {
		defaultFiller = newDefaultFiller()
	}

	return defaultFiller
}

func newDefaultFiller() *Filler {
	funcs := make(map[reflect.Kind]fillerFunc, 0)
	funcs[reflect.Bool] = func(field *fieldData, defaultValue string) {
		value, _ := strconv.ParseBool(defaultValue)
		field.Value.SetBool(value)
	}

	funcs[reflect.Int] = func(field *fieldData, defaultValue string) {
		value, _ := strconv.ParseInt(defaultValue, 10, 64)
		field.Value.SetInt(value)
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = funcs[reflect.Int]

	funcs[reflect.Float32] = func(field *fieldData, defaultValue string) {
		value, _ := strconv.ParseFloat(defaultValue, 64)
		field.Value.SetFloat(value)
	}

	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *fieldData, defaultValue string) {
		value, _ := strconv.ParseUint(defaultValue, 10, 64)
		field.Value.SetUint(value)
	}

	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.String] = func(field *fieldData, defaultValue string) {
		field.Value.SetString(defaultValue)
	}

	funcs[reflect.Slice] = func(field *fieldData, defaultValue string) {
		if field.Value.Type().Elem().Kind() == reflect.Uint8 {
			if field.Value.Bytes() != nil {
				return
			}

			field.Value.SetBytes([]byte(defaultValue))
		}
	}

	funcs[reflect.Struct] = func(field *fieldData, defaultValue string) {
		fields := getDefaultFiller().getFieldsFromValue(field.Value)
		getDefaultFiller().setDefaultValues(fields)
	}

	return &Filler{FuncByKind: funcs, Tag: "default"}
}
