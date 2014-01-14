package defaults

import (
	"reflect"
	"strconv"
)

type Field struct {
	Field reflect.StructField
	Value reflect.Value
}

func SetDefaults(variable interface{}) {
	fields := getFields(variable)
	setDefaultValues(fields)
}

func getFields(variable interface{}) []*Field {
	valueObject := reflect.ValueOf(variable).Elem()

	return getFieldsFromValue(valueObject)
}

func getFieldsFromValue(valueObject reflect.Value) []*Field {
	typeObject := valueObject.Type()

	count := valueObject.NumField()
	results := make([]*Field, 0)
	for i := 0; i < count; i++ {
		value := valueObject.Field(i)
		field := typeObject.Field(i)

		if value.CanSet() {
			results = append(results, &Field{
				Value: value,
				Field: field,
			})
		}
	}

	return results
}

func setDefaultValues(fields []*Field) {
	for _, field := range fields {
		setDefaultValue(field)
	}
}

func setDefaultValue(field *Field) {
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

func setDefaultValueToBool(field *Field, defaultValue string) {
	value, _ := strconv.ParseBool(defaultValue)
	field.Value.SetBool(value)
}

func setDefaultValueToInt(field *Field, defaultValue string) {
	value, _ := strconv.ParseInt(defaultValue, 10, 64)
	field.Value.SetInt(value)
}

func setDefaultValueToFloat(field *Field, defaultValue string) {
	value, _ := strconv.ParseFloat(defaultValue, 64)
	field.Value.SetFloat(value)
}

func setDefaultValueToUint(field *Field, defaultValue string) {
	value, _ := strconv.ParseUint(defaultValue, 10, 64)
	field.Value.SetUint(value)
}

func setDefaultValueToSlice(field *Field, defaultValue string) {
	if field.Value.Type().Elem().Kind() == reflect.Uint8 {
		field.Value.SetBytes([]byte(defaultValue))
	}
}

func setDefaultValueToString(field *Field, defaultValue string) {
	field.Value.SetString(defaultValue)
}

func setDefaultValueToStruct(field *Field, defaultValue string) {
	fields := getFieldsFromValue(field.Value)
	setDefaultValues(fields)
}
