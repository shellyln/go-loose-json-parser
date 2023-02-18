package marshal

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"
)

var (
	typeOfSliceOfAny  = reflect.TypeOf([]interface{}{})
	typeOfMapOfStrAny = reflect.TypeOf(map[string]interface{}{})
	typeOfTime        = reflect.TypeOf(time.Time{})
)

func unmarshalInterface(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Bool,
		reflect.String:
		rvTo.Set(rvFrom)

	case reflect.Slice, reflect.Array:
		matched := false
		if rvFrom.IsValid() {
			v := rvFrom.Interface()
			switch z := v.(type) {
			case []byte:
				rvTo.Set(reflect.ValueOf(base64.StdEncoding.EncodeToString(z)))
				matched = true
			}
		}
		if !matched {
			rvDest := reflect.New(typeOfSliceOfAny).Elem()
			unmarshalCore(rvFrom, rvDest, ctx, true)
			rvTo.Set(rvDest)
		}

	case reflect.Map, reflect.Struct:
		matched := false
		if rvFrom.IsValid() {
			v := rvFrom.Interface()
			switch z := v.(type) {
			case time.Time:
				if b, err := json.Marshal(z); err != nil {
					return err
				} else {
					s := string(b)
					rvTo.Set(reflect.ValueOf(s[1 : len(s)-1]))
					matched = true
				}
			}
		}
		if !matched {
			rvDest := reflect.New(typeOfMapOfStrAny).Elem()
			unmarshalCore(rvFrom, rvDest, ctx, true)
			rvTo.Set(rvDest)
		}

	default:
		// do nothing
	}

	return nil
}
