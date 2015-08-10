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

type Filler struct {
	FuncByType map[TypeHash]fillerFunc
	FuncByKind map[reflect.Kind]fillerFunc
	Tag        string
}

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
	results := make([]*fieldData, 0)
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

	function := f.getFunctionByType(field.Field.Type)
	if function != nil {
		function(field, tagValue)
		return
	}

	function = f.getFunctionByKind(field.Field.Type.Kind())
	if function != nil {
		function(field, tagValue)
		return
	}

	return
}

func (f *Filler) getFunctionByType(t reflect.Type) fillerFunc {
	if f, ok := f.FuncByType[GetTypeHash(t)]; ok == true {
		return f
	}

	return nil
}

func (f *Filler) getFunctionByKind(k reflect.Kind) fillerFunc {
	if f, ok := f.FuncByKind[k]; ok == true {
		return f
	}

	return nil
}

type TypeHash string

func GetTypeHash(t reflect.Type) TypeHash {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return TypeHash(fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()))
}
