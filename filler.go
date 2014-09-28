package defaults

import (
	"reflect"
)

type fieldData struct {
	Field reflect.StructField
	Value reflect.Value
}

type fillerFunc func(field *fieldData, config string)

type Filler struct {
	FuncByName map[string]fillerFunc
	FuncByKind map[reflect.Kind]fillerFunc
	Tag        string
}

func (self *Filler) Fill(variable interface{}) {
	fields := self.getFields(variable)
	self.setDefaultValues(fields)
}

func (self *Filler) getFields(variable interface{}) []*fieldData {
	valueObject := reflect.ValueOf(variable).Elem()

	return self.getFieldsFromValue(valueObject)
}

func (self *Filler) getFieldsFromValue(valueObject reflect.Value) []*fieldData {
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

func (self *Filler) setDefaultValues(fields []*fieldData) {
	for _, field := range fields {
		if self.isEmpty(field) {
			self.setDefaultValue(field)
		}
	}
}

func (self *Filler) isEmpty(field *fieldData) bool {
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
		if field.Value.Type().Elem().Kind() == reflect.Uint8 {
			if field.Value.Bytes() != nil {
				return false
			}
		}
	case reflect.String:
		if field.Value.String() != "" {
			return false
		}
	}

	return true
}

func (self *Filler) setDefaultValue(field *fieldData) {
	tagValue := field.Field.Tag.Get(self.Tag)

	function := self.getFunctionByKind(field.Field.Type.Kind())
	if function == nil {
		return
	}

	function(field, tagValue)
}

func (self *Filler) getFunctionByKind(k reflect.Kind) fillerFunc {
	if f, ok := self.FuncByKind[k]; ok == true {
		return f
	}

	return nil
}
