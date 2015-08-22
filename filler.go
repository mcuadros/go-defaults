package defaults

import (
	"fmt"
	"reflect"
)

type fieldData struct {
	Field reflect.StructField
	Value reflect.Value
}

type fillerFunc func(field *fieldData, config string)

// Filler contains all the functions to fill any struct field with any type
// allowing to define function by Kind, Type of field name
type Filler struct {
	FuncByName map[string]fillerFunc
	FuncByType map[TypeHash]fillerFunc
	FuncByKind map[reflect.Kind]fillerFunc
	Tag        string
}

// Fill apply all the functions contained on Filler, setting all the possible
// values
func (f *Filler) Fill(variable interface{}) {
	fields := f.getFields(variable)
	f.setDefaultValues(fields)
}

func (f *Filler) getFields(variable interface{}) []*fieldData {
	valueObject := reflect.ValueOf(variable).Elem()

	return f.getFieldsFromValue(valueObject)
}

func (f *Filler) getFieldsFromValue(valueObject reflect.Value) []*fieldData {
	typeObject := valueObject.Type()

	count := valueObject.NumField()
	var results []*fieldData
	for i := 0; i < count; i++ {
		value := valueObject.Field(i)
		field := typeObject.Field(i)

		if value.CanSet() {
			results = append(results, &fieldData{
				Value: value,
				Field: field,
			})
		}
	}

	return results
}

func (f *Filler) setDefaultValues(fields []*fieldData) {
	for _, field := range fields {
		if f.isEmpty(field) {
			f.setDefaultValue(field)
		}
	}
}

func (f *Filler) isEmpty(field *fieldData) bool {
	switch field.Value.Kind() {
	case reflect.Bool:
		if field.Value.Bool() != false {
			return false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Value.Int() != 0 {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if field.Value.Float() != .0 {
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Value.Uint() != 0 {
			return false
		}
	case reflect.Slice:
		if field.Value.Len() != 0 {
			return false
		}
	case reflect.String:
		if field.Value.String() != "" {
			return false
		}
	}

	return true
}

func (f *Filler) setDefaultValue(field *fieldData) {
	tagValue := field.Field.Tag.Get(f.Tag)

	getters := []func(field *fieldData) fillerFunc{
		f.getFunctionByName,
		f.getFunctionByType,
		f.getFunctionByKind,
	}

	for _, getter := range getters {
		filler := getter(field)
		if filler != nil {
			filler(field, tagValue)
			return
		}
	}

	return
}

func (f *Filler) getFunctionByName(field *fieldData) fillerFunc {
	if f, ok := f.FuncByName[field.Field.Name]; ok == true {
		return f
	}

	return nil
}

func (f *Filler) getFunctionByType(field *fieldData) fillerFunc {
	if f, ok := f.FuncByType[GetTypeHash(field.Field.Type)]; ok == true {
		return f
	}

	return nil
}

func (f *Filler) getFunctionByKind(field *fieldData) fillerFunc {
	if f, ok := f.FuncByKind[field.Field.Type.Kind()]; ok == true {
		return f
	}

	return nil
}

// TypeHash is a string representing a reflect.Type following the next pattern:
// <package.name>.<type.name>
type TypeHash string

// GetTypeHash returns the TypeHash for a given reflect.Type
func GetTypeHash(t reflect.Type) TypeHash {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return TypeHash(fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()))
}
