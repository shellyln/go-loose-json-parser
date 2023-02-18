package marshal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"strconv"
	"time"
)

func unmarshalInt(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvTo.SetInt(rvFrom.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvTo.SetInt(int64(rvFrom.Uint()))
	case reflect.Float32, reflect.Float64:
		rvTo.SetInt(int64(rvFrom.Float()))
	case reflect.Complex64, reflect.Complex128:
		rvTo.SetInt(int64(real(rvFrom.Complex())))
	case reflect.Bool:
		v := rvFrom.Bool()
		if v {
			rvTo.SetInt(1)
		} else {
			rvTo.SetInt(0)
		}
	case reflect.String:
		v := rvFrom.String()
		if z, err := strconv.ParseInt(v, 10, 64); err != nil {
			return err
		} else {
			rvTo.SetInt(z)
		}
	default:
		return fmt.Errorf("Type unmatched: %v -> Int", rvFrom.Interface())
	}

	return nil
}

func unmarshalUint(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvTo.SetUint(uint64(rvFrom.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvTo.SetUint(rvFrom.Uint())
	case reflect.Float32, reflect.Float64:
		rvTo.SetUint(uint64(rvFrom.Float()))
	case reflect.Complex64, reflect.Complex128:
		rvTo.SetUint(uint64(real(rvFrom.Complex())))
	case reflect.Bool:
		v := rvFrom.Bool()
		if v {
			rvTo.SetUint(1)
		} else {
			rvTo.SetUint(0)
		}
	case reflect.String:
		v := rvFrom.String()
		if z, err := strconv.ParseUint(v, 10, 64); err != nil {
			return err
		} else {
			rvTo.SetUint(z)
		}
	default:
		return fmt.Errorf("Type unmatched: %v -> Uint", rvFrom.Interface())
	}

	return nil
}

func setNanOrInfMap(rvFrom, rvTo reflect.Value, ctx *marshalContext) bool {
	matched := false
	if rvFrom.Type().Key().Kind() == reflect.String {
		if !matched {
			rvNan := rvFrom.MapIndex(reflect.ValueOf("nan"))
			for rvNan.IsValid() {
				switch rvNan.Kind() {
				case reflect.Pointer, reflect.Interface:
					rvNan = rvNan.Elem()
					continue
				case reflect.Bool:
					rvTo.SetFloat(math.NaN())
					matched = true
				}
				break
			}
		}
		if !matched {
			rvInf := rvFrom.MapIndex(reflect.ValueOf("inf"))
			for rvInf.IsValid() {
				switch rvInf.Kind() {
				case reflect.Pointer, reflect.Interface:
					rvInf = rvInf.Elem()
					continue
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					sign := 1
					if z := rvInf.Int(); z < 0 {
						sign = -1
					}
					rvTo.SetFloat(math.Inf(sign))
					matched = true
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					rvTo.SetFloat(math.Inf(1))
					matched = true
				case reflect.Float32, reflect.Float64:
					sign := 1
					if z := rvInf.Float(); z < 0 {
						sign = -1
					}
					rvTo.SetFloat(math.Inf(sign))
					matched = true
				}
				break
			}
		}
	}
	return matched
}

func unmarshalFloat(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvTo.SetFloat(float64(rvFrom.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvTo.SetFloat(float64(rvFrom.Uint()))
	case reflect.Float32, reflect.Float64:
		rvTo.SetFloat(rvFrom.Float())
	case reflect.Complex64, reflect.Complex128:
		rvTo.SetFloat(real(rvFrom.Complex()))
	case reflect.Bool:
		v := rvFrom.Bool()
		if v {
			rvTo.SetFloat(1)
		} else {
			rvTo.SetFloat(0)
		}
	case reflect.String:
		v := rvFrom.String()
		if z, err := strconv.ParseFloat(v, 64); err != nil {
			return err
		} else {
			rvTo.SetFloat(z)
		}
	case reflect.Map:
		if !setNanOrInfMap(rvFrom, rvTo, ctx) {
			return fmt.Errorf("Type unmatched: %v -> Float", rvFrom.Interface())
		}
	default:
		return fmt.Errorf("Type unmatched: %v -> Float", rvFrom.Interface())
	}

	return nil
}

func setComplexMap(rvFrom, rvTo reflect.Value, ctx *marshalContext) bool {
	matched := false
	var re, im float64
	if rvFrom.Type().Key().Kind() == reflect.String {
		rvRe := rvFrom.MapIndex(reflect.ValueOf("re"))
		for rvRe.IsValid() {
			switch rvRe.Kind() {
			case reflect.Pointer, reflect.Interface:
				rvRe = rvRe.Elem()
				continue
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				re = float64(rvRe.Int())
				matched = true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				re = float64(rvRe.Uint())
				matched = true
			case reflect.Float32, reflect.Float64:
				re = rvRe.Float()
				matched = true
			case reflect.Map:
				s := make([]float64, 1)
				matched = setNanOrInfMap(rvRe, reflect.ValueOf(s).Index(0), ctx)
				re = s[0]
			}
			break
		}
		rvIm := rvFrom.MapIndex(reflect.ValueOf("im"))
		for rvIm.IsValid() {
			switch rvIm.Kind() {
			case reflect.Pointer, reflect.Interface:
				rvIm = rvIm.Elem()
				continue
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				im = float64(rvIm.Int())
				matched = true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				im = float64(rvIm.Uint())
				matched = true
			case reflect.Float32, reflect.Float64:
				im = rvIm.Float()
				matched = true
			case reflect.Map:
				s := make([]float64, 1)
				matched = setNanOrInfMap(rvIm, reflect.ValueOf(s).Index(0), ctx)
				im = s[0]
			}
			break
		}
		if matched {
			rvTo.SetComplex(complex(re, im))
		}
	}
	return matched
}

func unmarshalComplex(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvTo.SetComplex(complex(float64(rvFrom.Int()), 0))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvTo.SetComplex(complex(float64(rvFrom.Uint()), 0))
	case reflect.Float32, reflect.Float64:
		rvTo.SetComplex(complex(rvFrom.Float(), 0))
	case reflect.Complex64, reflect.Complex128:
		rvTo.SetComplex(rvFrom.Complex())
	case reflect.Bool:
		v := rvFrom.Bool()
		if v {
			rvTo.SetComplex(1)
		} else {
			rvTo.SetComplex(0)
		}
	case reflect.String:
		v := rvFrom.String()
		if z, err := strconv.ParseComplex(v, 64); err != nil {
			return err
		} else {
			rvTo.SetComplex(z)
		}
	case reflect.Map:
		if !setComplexMap(rvFrom, rvTo, ctx) {
			return fmt.Errorf("Type unmatched: %v -> Complex", rvFrom.Interface())
		}
	default:
		return fmt.Errorf("Type unmatched: %v -> Complex", rvFrom.Interface())
	}

	return nil
}

func unmarshalBool(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v := rvFrom.Int(); v != 0 {
			rvTo.SetBool(true)
		} else {
			rvTo.SetBool(false)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v := rvFrom.Uint(); v != 0 {
			rvTo.SetBool(true)
		} else {
			rvTo.SetBool(false)
		}
	case reflect.Float32, reflect.Float64:
		if v := rvFrom.Float(); !math.IsNaN(v) && v != 0 {
			rvTo.SetBool(true)
		} else {
			rvTo.SetBool(false)
		}
	case reflect.Complex64, reflect.Complex128:
		if v := rvFrom.Complex(); !cmplx.IsNaN(v) && v != 0 {
			rvTo.SetBool(true)
		} else {
			rvTo.SetBool(false)
		}
	case reflect.Bool:
		rvTo.SetBool(rvFrom.Bool())
	case reflect.String:
		v := rvFrom.String()
		if z, err := strconv.ParseBool(v); err != nil {
			return err
		} else {
			rvTo.SetBool(z)
		}
	default:
		return fmt.Errorf("Type unmatched: %v -> Bool", rvFrom.Interface())
	}

	return nil
}

func unmarshalString(rvFrom, rvTo reflect.Value, ctx *marshalContext) error {
	switch rvFrom.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvTo.SetString(strconv.FormatInt(rvFrom.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvTo.SetString(strconv.FormatUint(rvFrom.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		rvTo.SetString(strconv.FormatFloat(rvFrom.Float(), 'g', -1, 64))
	case reflect.Complex64, reflect.Complex128:
		s := strconv.FormatComplex(rvFrom.Complex(), 'g', -1, 128)
		rvTo.SetString(s[1 : len(s)-1])
	case reflect.Bool:
		v := rvFrom.Bool()
		if v {
			rvTo.SetString("true")
		} else {
			rvTo.SetString("false")
		}
	case reflect.String:
		rvTo.SetString(rvFrom.String())
	default:
		matched := false
		if rvFrom.IsValid() {
			v := rvFrom.Interface()
			switch z := v.(type) {
			case []byte:
				rvTo.SetString(base64.StdEncoding.EncodeToString(z))
				matched = true
			case time.Time:
				if b, err := json.Marshal(z); err != nil {
					return err
				} else {
					s := string(b)
					rvTo.SetString(s[1 : len(s)-1])
					matched = true
				}
			}
		}
		if !matched {
			rvTo.SetString(fmt.Sprintf("%v", rvFrom.Interface()))
		}
	}

	return nil
}
