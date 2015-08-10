package defaults

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"
)

func Factory(variable interface{}) {
	getFactoryFiller().Fill(variable)
}

var factoryFiller *Filler = nil

func getFactoryFiller() *Filler {
	if factoryFiller == nil {
		factoryFiller = newFactoryFiller()
	}

	return factoryFiller
}

func newFactoryFiller() *Filler {
	rand.Seed(time.Now().UTC().UnixNano())

	funcs := make(map[reflect.Kind]fillerFunc, 0)

	funcs[reflect.Bool] = func(field *fieldData, _ string) {
		if rand.Intn(1) == 1 {
			field.Value.SetBool(true)
		} else {
			field.Value.SetBool(false)
		}
	}

	funcs[reflect.Int] = func(field *fieldData, _ string) {
		field.Value.SetInt(int64(rand.Int()))
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = funcs[reflect.Int]

	funcs[reflect.Float32] = func(field *fieldData, tagValue string) {
		field.Value.SetFloat(rand.Float64())
	}

	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *fieldData, tagValue string) {
		field.Value.SetUint(uint64(rand.Uint32()))
	}

	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.String] = func(field *fieldData, tagValue string) {
		field.Value.SetString(randomString())
	}

	funcs[reflect.Slice] = func(field *fieldData, tagValue string) {
		if field.Value.Type().Elem().Kind() == reflect.Uint8 {
			if field.Value.Bytes() != nil {
				return
			}

			field.Value.SetBytes([]byte(randomString()))
		}
	}

	funcs[reflect.Struct] = func(field *fieldData, tagValue string) {
		fields := getDefaultFiller().getFieldsFromValue(field.Value)
		getDefaultFiller().setDefaultValues(fields)
	}

	return &Filler{FuncByKind: funcs, Tag: "factory"}
}

func randomString() string {
	hash := md5.Sum([]byte(time.Now().UTC().String()))
	return hex.EncodeToString(hash[:])
}
