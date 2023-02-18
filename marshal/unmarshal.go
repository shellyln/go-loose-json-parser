package marshal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

type IMarshal interface {
	MarshalLp() (to interface{}, err error)
}

type IUnmarshal interface {
	UnmarshalLp(from interface{}) error
}

type MarshalOptions struct {
	TagName                string // Tag name of the struct fields
	NoCopyUnexportedFields bool   // If true, no shallow copying of unexported fields
	NoCustomMarshaller     bool   // If true, IMarshal and IUnmarshal are not used
}

type marshalContext struct {
	opts MarshalOptions
	ptrs map[unsafe.Pointer]int
}

var marshalOptsDefault = MarshalOptions{
	TagName:                "json",
	NoCopyUnexportedFields: false,
}

func unmarshalCore(rvFrom, rvTo reflect.Value, ctx *marshalContext, rvFromReused bool) error {
	rtTo := rvTo.Type()

	if !rvFromReused {
		switch rvFrom.Kind() {
		case reflect.Pointer, reflect.Map, reflect.Slice:
			ptr := rvFrom.UnsafePointer()
			if _, ok := ctx.ptrs[ptr]; ok {
				// circular references
				return nil
			}
			ctx.ptrs[ptr] = 0
			defer delete(ctx.ptrs, ptr)
		}
	}

	if !rvTo.CanSet() {
		// it is unexported or unaddressable field
		return nil
	}
	if !rvFrom.CanInterface() {
		// it is unexported field
		return nil
	}

	if !ctx.opts.NoCustomMarshaller && !rvFromReused {
		if imar, ok := rvFrom.Interface().(IMarshal); ok {
			if tmp, err := imar.MarshalLp(); err != nil {
				return err
			} else {
				// NOTE: Treat as a reuse of rvFrom
				return unmarshalCore(reflect.ValueOf(tmp), rvTo, ctx, true)
			}
		} else {
			ptrFrom := reflect.New(rvFrom.Type())
			if imar, ok := ptrFrom.Interface().(IMarshal); ok {
				ptrFrom.Elem().Set(rvFrom)
				if tmp, err := imar.MarshalLp(); err != nil {
					return err
				} else {
					// NOTE: Treat as a reuse of rvFrom
					return unmarshalCore(reflect.ValueOf(tmp), rvTo, ctx, true)
				}
			}
		}
	}

	switch rvFrom.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !rvFrom.IsNil() {
			return unmarshalCore(rvFrom.Elem(), rvTo, ctx, false)
		} else {
			return nil
		}
	case reflect.Slice, reflect.Map:
		if rvFrom.IsNil() {
			return nil
		}
	}

	if !ctx.opts.NoCustomMarshaller {
		if iumar, ok := rvTo.Interface().(IUnmarshal); ok {
			if err := iumar.UnmarshalLp(rvFrom.Interface()); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			ptrTo := rvTo.Addr()
			if iumar, ok := ptrTo.Interface().(IUnmarshal); ok {
				if err := iumar.UnmarshalLp(rvFrom.Interface()); err != nil {
					return err
				} else {
					return nil
				}
			}
		}
	}

	switch rvTo.Kind() {
	case reflect.Pointer:
		reflect.New(rtTo.Elem())
		rvTo.Set(reflect.New(rtTo.Elem()))
		return unmarshalCore(rvFrom, rvTo.Elem(), ctx, true)

	case reflect.Interface:
		return unmarshalInterface(rvFrom, rvTo, ctx)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return unmarshalInt(rvFrom, rvTo, ctx)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return unmarshalUint(rvFrom, rvTo, ctx)
	case reflect.Float32, reflect.Float64:
		return unmarshalFloat(rvFrom, rvTo, ctx)
	case reflect.Complex64, reflect.Complex128:
		return unmarshalComplex(rvFrom, rvTo, ctx)
	case reflect.Bool:
		return unmarshalBool(rvFrom, rvTo, ctx)
	case reflect.String:
		return unmarshalString(rvFrom, rvTo, ctx)

	case reflect.Slice, reflect.Array:
		switch rvFrom.Kind() {
		case reflect.Slice, reflect.Array:
			length := rvFrom.Len()
			if rvTo.Kind() == reflect.Slice {
				rvTo.Set(reflect.MakeSlice(reflect.SliceOf(rtTo.Elem()), length, length))
			} else {
				if rvTo.Len() < length {
					length = rvTo.Len()
				}
			}
			for i := 0; i < length; i++ {
				if err := unmarshalCore(rvFrom.Index(i), rvTo.Index(i), ctx, false); err != nil {
					return err
				}
			}

		case reflect.String:
			if rtTo.Kind() == reflect.Slice && rtTo.Elem() == reflect.TypeOf(byte(0)) {
				if dst, err := base64.StdEncoding.DecodeString(rvFrom.String()); err != nil {
					return err
				} else {
					rvTo.Set(reflect.ValueOf(dst))
				}
			} else {
				return fmt.Errorf("Type unmatched: %v -> Slice or Array", rvFrom.Interface())
			}

		default:
			return fmt.Errorf("Type unmatched: %v -> Slice or Array", rvFrom.Interface())
		}

	case reflect.Map:
		rvTo.Set(reflect.MakeMap(rtTo))

		switch rvFrom.Kind() {
		case reflect.Map:
			if !rvFrom.Type().Key().AssignableTo(rtTo.Key()) {
				return fmt.Errorf("Map -> map: Key types are incompatible")
			}

			for _, rvSrcKey := range rvFrom.MapKeys() {
				rvSrcValue := rvFrom.MapIndex(rvSrcKey)
				rvDest := reflect.New(rtTo.Elem()).Elem()
				if err := unmarshalCore(rvSrcValue, rvDest, ctx, false); err != nil {
					return err
				}
				rvTo.SetMapIndex(rvSrcKey, rvDest)
			}

		case reflect.Struct:
			if rtTo.Key().Kind() != reflect.String {
				return fmt.Errorf("Struct -> map: Map key type should be string (parameter \"to\")")
			}

			length := rvFrom.NumField()
			hasTagName := false
			if ctx.opts.TagName != "" {
				hasTagName = true
			}
			for i := 0; i < length; i++ {
				rtSrcField := rvFrom.Type().Field(i)
				srcFieldName := rtSrcField.Name
				if hasTagName {
					tags := strings.Split(rtSrcField.Tag.Get(ctx.opts.TagName), ",")
					if tags[0] != "" {
						srcFieldName = tags[0]
					}
				}
				rvDest := reflect.New(rtTo.Elem()).Elem()
				if err := unmarshalCore(rvFrom.Field(i), rvDest, ctx, false); err != nil {
					return err
				}
				rvTo.SetMapIndex(reflect.ValueOf(srcFieldName), rvDest)
			}

		default:
			return fmt.Errorf("Type unmatched: %v -> Map", rvFrom.Interface())
		}

	case reflect.Struct:
		switch rvFrom.Kind() {
		case reflect.Struct:
			length := rvTo.NumField()
			if rtTo == rvFrom.Type() {
				// shallow copy all the fields
				if !ctx.opts.NoCopyUnexportedFields {
					rvTo.Set(rvFrom)
				}
			}
			// deep copy
			for i := 0; i < length; i++ {
				rtDestField := rtTo.Field(i)
				destFieldName := rtDestField.Name
				if err := unmarshalCore(rvFrom.FieldByName(destFieldName), rvTo.Field(i), ctx, false); err != nil {
					return err
				}
			}

		case reflect.Map:
			if rvFrom.Type().Key().Kind() != reflect.String {
				return fmt.Errorf("Map -> struct: Map key type should be string (parameter \"from\")")
			}
			length := rvTo.NumField()
			hasTagName := false
			if ctx.opts.TagName != "" {
				hasTagName = true
			}
			for i := 0; i < length; i++ {
				rtDestField := rtTo.Field(i)
				destFieldName := rtDestField.Name
				if hasTagName {
					destTags := strings.Split(rtDestField.Tag.Get(ctx.opts.TagName), ",")
					if destTags[0] != "" {
						destFieldName = destTags[0]
					}
				}
				rvSrcValue := rvFrom.MapIndex(reflect.ValueOf(destFieldName))
				if rvSrcValue.IsValid() {
					if err := unmarshalCore(rvSrcValue, rvTo.Field(i), ctx, false); err != nil {
						return err
					}
				}
			}

		case reflect.String:
			matched := false
			if rtTo == typeOfTime {
				var v time.Time
				if err := json.Unmarshal([]byte("\""+rvFrom.String()+"\""), &v); err == nil {
					rvTo.Set(reflect.ValueOf(v))
					matched = true
				}
			}
			if !matched {
				return fmt.Errorf("Type unmatched: %v -> Struct", rvFrom.Interface())
			}

		default:
			return fmt.Errorf("Type unmatched: %v -> Struct", rvFrom.Interface())
		}

	default:
	}

	return nil
}

func Unmarshal(from interface{}, to interface{}, opts *MarshalOptions) error {
	rvTo := reflect.ValueOf(to)
	if rvTo.Kind() != reflect.Pointer {
		return fmt.Errorf("2nd Parameter \"to\" should be pointer")
	}

	options := opts
	if options == nil {
		options = &marshalOptsDefault
	}
	ctx := &marshalContext{
		opts: *options,
		ptrs: make(map[unsafe.Pointer]int),
	}

	return unmarshalCore(reflect.ValueOf(from), rvTo.Elem(), ctx, false)
}
