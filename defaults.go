package defaults

import (
	"reflect"
	"strconv"
)

type fieldData struct {
	Field reflect.StructField
	Value reflect.Value
}

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
	fields := getFields(variable)
	setDefaultValues(fields)
}

func getFields(variable interface{}) []*fieldData {
	valueObject := reflect.ValueOf(variable).Elem()

	return getFieldsFromValue(valueObject)
}

func getFieldsFromValue(valueObject reflect.Value) []*fieldData {
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

func setDefaultValues(fields []*fieldData) {
	for _, field := range fields {
		setDefaultValue(field)
	}
}

func setDefaultValue(field *fieldData) {
	defaultValue := field.Field.Tag.Get("default")

	switch field.Value.Kind() {
	case reflect.Bool:
		setDefaultValueToBool(field, defaultValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setDefaultValueToInt(field, defaultValue)
	case reflect.Float32, reflect.Float64:
		setDefaultValueToFloat(field, defaultValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setDefaultValueToUint(field, defaultValue)
	case reflect.Slice:
		setDefaultValueToSlice(field, defaultValue)
	case reflect.String:
		setDefaultValueToString(field, defaultValue)
	case reflect.Struct:
		setDefaultValueToStruct(field, defaultValue)
	}
}

func setDefaultValueToBool(field *fieldData, defaultValue string) {
	if field.Value.Bool() != false {
		return
	}

	value, _ := strconv.ParseBool(defaultValue)
	field.Value.SetBool(value)
}

func setDefaultValueToInt(field *fieldData, defaultValue string) {
	if field.Value.Int() != 0 {
		return
	}

	value, _ := strconv.ParseInt(defaultValue, 10, 64)
	field.Value.SetInt(value)
}

func setDefaultValueToFloat(field *fieldData, defaultValue string) {
	if field.Value.Float() != .0 {
		return
	}

	value, _ := strconv.ParseFloat(defaultValue, 64)
	field.Value.SetFloat(value)
}

func setDefaultValueToUint(field *fieldData, defaultValue string) {
	if field.Value.Uint() != 0 {
		return
	}

	value, _ := strconv.ParseUint(defaultValue, 10, 64)
	field.Value.SetUint(value)
}

func setDefaultValueToSlice(field *fieldData, defaultValue string) {
	if field.Value.Type().Elem().Kind() == reflect.Uint8 {
		if field.Value.Bytes() != nil {
			return
		}

		field.Value.SetBytes([]byte(defaultValue))
	}
}

func setDefaultValueToString(field *fieldData, defaultValue string) {
	if field.Value.String() != "" {
		return
	}

	field.Value.SetString(defaultValue)
}

func setDefaultValueToStruct(field *fieldData, defaultValue string) {
	fields := getFieldsFromValue(field.Value)
	setDefaultValues(fields)
}
