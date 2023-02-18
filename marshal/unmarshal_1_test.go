package marshal_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/shellyln/go-loose-json-parser/marshal"
)

func TestPrimitive1(t *testing.T) {
	src := int(1234)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1b(t *testing.T) {
	src := uint(1234)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1c(t *testing.T) {
	src := float64(1234)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1d(t *testing.T) {
	src := complex(1234, 5678)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1e(t *testing.T) {
	src := true

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1f(t *testing.T) {
	src := false

	dst := int(1)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(0)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1g(t *testing.T) {
	src := "1234"

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive1g2(t *testing.T) {
	src := "q1234"

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func TestPrimitive2(t *testing.T) {
	src := int8(123)

	dst := int8(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int8(123)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive3(t *testing.T) {
	src := int16(1234)

	dst := int16(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int16(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive4(t *testing.T) {
	src := int32(1234)

	dst := int32(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int32(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive5(t *testing.T) {
	src := int64(1234)

	dst := int64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6(t *testing.T) {
	src := uint(1234)

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6b(t *testing.T) {
	src := int(1234)

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6c(t *testing.T) {
	src := float64(1234)

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6d(t *testing.T) {
	src := complex(1234, 5678)

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6e(t *testing.T) {
	src := true

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6f(t *testing.T) {
	src := false

	dst := uint(1234)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(0)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6g(t *testing.T) {
	src := "1234"

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive6g2(t *testing.T) {
	src := "q1234"

	dst := uint(0)
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func TestPrimitive7(t *testing.T) {
	src := uint8(123)

	dst := uint8(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint8(123)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive8(t *testing.T) {
	src := uint16(1234)

	dst := uint16(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint16(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive9(t *testing.T) {
	src := uint32(1234)

	dst := uint32(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint32(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive10(t *testing.T) {
	src := uint64(1234)

	dst := uint64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uint64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive11(t *testing.T) {
	src := uintptr(1234)

	dst := uintptr(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := uintptr(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive12(t *testing.T) {
	src := float32(1234)

	dst := float32(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float32(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13(t *testing.T) {
	src := float64(1234)

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13b(t *testing.T) {
	src := int(1234)

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13c(t *testing.T) {
	src := uint(1234)

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13d(t *testing.T) {
	src := complex(1234, 5678)

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13e(t *testing.T) {
	src := true

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13f(t *testing.T) {
	src := false

	dst := float64(1234)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(0)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13g(t *testing.T) {
	src := "1234"

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := float64(1234)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive13g2(t *testing.T) {
	src := "q1234"

	dst := float64(0)
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func TestPrimitive14(t *testing.T) {
	src := complex64(complex(1234, 5678))

	dst := complex64(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex64(complex(1234, 5678))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15(t *testing.T) {
	src := complex128(complex(1234, 5678))

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(1234, 5678))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15b(t *testing.T) {
	src := int(1234)

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(1234, 0))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15c(t *testing.T) {
	src := uint(1234)

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(1234, 0))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15d(t *testing.T) {
	src := float64(1234)

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(1234, 0))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15e(t *testing.T) {
	src := true

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(1, 0))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15f(t *testing.T) {
	src := false

	dst := complex128(complex(1234, 5678))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(0, 0))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15g(t *testing.T) {
	src := "-1234+5678i"

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := complex128(complex(-1234, 5678))
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive15g2(t *testing.T) {
	src := "q1234"

	dst := complex128(complex(0, 0))
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func TestPrimitive16(t *testing.T) {
	src := true

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16b(t *testing.T) {
	src := int(1)

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16c(t *testing.T) {
	src := int(0)

	dst := true
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := false
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16d(t *testing.T) {
	src := uint(1)

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16e(t *testing.T) {
	src := uint(0)

	dst := true
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := false
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16f(t *testing.T) {
	src := float64(1)

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16g(t *testing.T) {
	src := float64(0)

	dst := true
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := false
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16h(t *testing.T) {
	src := complex(0, 1)

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16i(t *testing.T) {
	src := complex(0, 0)

	dst := true
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := false
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16j(t *testing.T) {
	src := "true"

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := true
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16k(t *testing.T) {
	src := "false"

	dst := true
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := false
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive16k2(t *testing.T) {
	src := "q1234"

	dst := false
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func TestPrimitive17(t *testing.T) {
	src := "1234"

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "1234"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17b(t *testing.T) {
	src := int(1234)

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "1234"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17c(t *testing.T) {
	src := uint(1234)

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "1234"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17d(t *testing.T) {
	src := float64(1234)

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "1234"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17e(t *testing.T) {
	src := complex(-1234, 5678)

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "-1234+5678i"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17e2(t *testing.T) {
	src := complex(math.NaN(), math.NaN())

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "NaN+NaNi"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17e3(t *testing.T) {
	src := complex(math.Inf(-1), math.Inf(-1))

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "-Inf-Infi"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17f(t *testing.T) {
	src := true

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "true"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestPrimitive17g(t *testing.T) {
	src := false

	dst := ""
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := "false"
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestArray1(t *testing.T) {
	src := []interface{}{1.0, 2.0, 3.0}

	dst := []int{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}
