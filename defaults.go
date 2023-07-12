package defaults

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Applies the default values to the struct object, the struct type must have
// the StructTag with name "default" and the directed value.
//
// Usage
//
//	type ExampleBasic struct {
//	    Foo bool   `default:"true"`
//	    Bar string `default:"33"`
//	    Qux int8
//	    Dur time.Duration `default:"2m3s"`
//	}
//
//	 foo := &ExampleBasic{}
//	 SetDefaults(foo)
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
	funcs := make(map[reflect.Kind]FillerFunc, 0)
	funcs[reflect.Bool] = func(field *FieldData) {
		value, _ := strconv.ParseBool(field.TagValue)
		field.Value.SetBool(value)
	}

	funcs[reflect.Int] = func(field *FieldData) {
		value, _ := strconv.ParseInt(field.TagValue, 10, 64)
		field.Value.SetInt(value)
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = func(field *FieldData) {
		if field.Field.Type == reflect.TypeOf(time.Second) {
			value, _ := time.ParseDuration(field.TagValue)
			field.Value.Set(reflect.ValueOf(value))
		} else {
			value, _ := strconv.ParseInt(field.TagValue, 10, 64)
			field.Value.SetInt(value)
		}
	}

	funcs[reflect.Float32] = func(field *FieldData) {
		value, _ := strconv.ParseFloat(field.TagValue, 64)
		field.Value.SetFloat(value)
	}

	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *FieldData) {
		value, _ := strconv.ParseUint(field.TagValue, 10, 64)
		field.Value.SetUint(value)
	}

	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.String] = func(field *FieldData) {
		tagValue := parseDateTimeString(field.TagValue)
		field.Value.SetString(tagValue)
	}

	funcs[reflect.Struct] = func(field *FieldData) {
		fields := getDefaultFiller().GetFieldsFromValue(field.Value, nil)
		getDefaultFiller().SetDefaultValues(fields)
	}

	types := make(map[TypeHash]FillerFunc, 1)
	types["time.Duration"] = func(field *FieldData) {
		d, _ := time.ParseDuration(field.TagValue)
		if field.Value.Kind() == reflect.Ptr {
			field.Value.Set(reflect.ValueOf(&d))
		} else {
			field.Value.Set(reflect.ValueOf(d))
		}
	}

	funcs[reflect.Slice] = func(field *FieldData) {
		k := field.Value.Type().Elem().Kind()
		switch k {
		case reflect.Uint8:
			if field.Value.Bytes() != nil {
				return
			}
			field.Value.SetBytes([]byte(field.TagValue))
		case reflect.Struct:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				fields := getDefaultFiller().GetFieldsFromValue(field.Value.Index(i), nil)
				getDefaultFiller().SetDefaultValues(fields)
			}
		case reflect.Ptr:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				if field.Value.Index(i).IsZero() {
					newValue := reflect.New(field.Value.Index(i).Type().Elem())
					field.Value.Index(i).Set(newValue)
				}
				fields := getDefaultFiller().GetFieldsFromValue(field.Value.Index(i).Elem(), nil)
				getDefaultFiller().SetDefaultValues(fields)
			}
		default:
			//处理形如 [1,2,3,4]
			reg := regexp.MustCompile(`^\[(.*)\]$`)
			matchs := reg.FindStringSubmatch(field.TagValue)
			if len(matchs) != 2 {
				return
			}
			if matchs[1] == "" {
				field.Value.Set(reflect.MakeSlice(field.Value.Type(), 0, 0))
			} else {
				defaultValue := strings.Split(matchs[1], ",")
				result := reflect.MakeSlice(field.Value.Type(), len(defaultValue), len(defaultValue))
				for i := 0; i < len(defaultValue); i++ {
					itemValue := result.Index(i)
					item := &FieldData{
						Value:    itemValue,
						Field:    reflect.StructField{},
						TagValue: defaultValue[i],
						Parent:   nil,
					}
					funcs[k](item)
				}
				field.Value.Set(result)
			}
		}
	}

	funcs[reflect.Ptr] = func(field *FieldData) {
		k := field.Value.Type().Elem().Kind()
		if k != reflect.Struct && k != reflect.Slice && k != reflect.Ptr && field.TagValue == "" {
			return
		}
		if field.Value.IsNil() {
			v := reflect.New(field.Value.Type().Elem())
			field.Value.Set(v)
		}
		elemField := &FieldData{
			Value: field.Value.Elem(),
			Field: reflect.StructField{
				Type: field.Field.Type.Elem(),
				Tag:  field.Field.Tag,
			},
			TagValue: field.TagValue,
			Parent:   nil,
		}
		funcs[field.Value.Elem().Kind()](elemField)
	}

	return &Filler{FuncByKind: funcs, FuncByType: types, Tag: "default"}
}

func parseDateTimeString(data string) string {

	pattern := regexp.MustCompile(`\{\{(\w+\:(?:-|)\d*,(?:-|)\d*,(?:-|)\d*)\}\}`)
	matches := pattern.FindAllStringSubmatch(data, -1) // matches is [][]string
	for _, match := range matches {

		tags := strings.Split(match[1], ":")
		if len(tags) == 2 {

			valueStrings := strings.Split(tags[1], ",")
			if len(valueStrings) == 3 {
				var values [3]int
				for key, valueString := range valueStrings {
					num, _ := strconv.ParseInt(valueString, 10, 64)
					values[key] = int(num)
				}

				switch tags[0] {

				case "date":
					str := time.Now().AddDate(values[0], values[1], values[2]).Format("2006-01-02")
					data = strings.Replace(data, match[0], str, -1)
				case "time":
					str := time.Now().Add((time.Duration(values[0]) * time.Hour) +
						(time.Duration(values[1]) * time.Minute) +
						(time.Duration(values[2]) * time.Second)).Format("15:04:05")
					data = strings.Replace(data, match[0], str, -1)
				}
			}
		}

	}
	return data
}
