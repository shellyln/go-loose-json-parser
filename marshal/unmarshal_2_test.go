package marshal_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/shellyln/go-loose-json-parser/marshal"
)

// Untyped -> Untyped
func Test1(t *testing.T) {
	src := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4, nil},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4, nil},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Untyped (+ typed array) -> Untyped
func Test1b(t *testing.T) {
	src := map[string]interface{}{
		"foo": 1,
		"bar": []int{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test1c(t *testing.T) {
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")
	src := map[string]interface{}{
		"foo": wt1,
		"bar": []int{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": "1999-12-31T23:45:59.123456789Z",
		"bar": []interface{}{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Untyped -> Typed
func Test2(t *testing.T) {
	type t02a struct {
		Aaa int `json:"aaa"`
	}
	type t02 struct {
		Foo int   `json:"foo"`
		Bar []int `json:"bar"`
		Baz t02a  `json:"baz"`
	}

	src := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t02{
		Foo: 1,
		Bar: []int{2, 3, 4},
		Baz: t02a{Aaa: 5},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Untyped (mixed type) -> Typed (fixed array)
func Test2b(t *testing.T) {
	type t02a struct {
		Aaa int `json:"aaa"`
	}
	type t02 struct {
		Foo int    `json:"foo"`
		Bar [2]int `json:"bar"`
		Baz t02a   `json:"baz"`
	}

	src := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{"2", float64(3), 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t02{
		Foo: 1,
		Bar: [2]int{2, 3},
		Baz: t02a{Aaa: 5},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Untyped (mixed type) -> Typed
func Test2c(t *testing.T) {
	type t02a struct {
		Aaa int `json:"aaa"`
	}
	type t02 struct {
		Foo int   `json:"foo"`
		Bar []int `json:"bar"`
		Baz t02a  `json:"baz"`
	}

	src := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{float64(2), 3, "4", nil},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t02{
		Foo: 1,
		Bar: []int{2, 3, 4, 0},
		Baz: t02a{Aaa: 5},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Typed -> Typed
func Test3(t *testing.T) {
	type t03a struct {
		Aaa int `json:"aaa"`
	}
	type t03 struct {
		Foo int   `json:"foo"`
		Bar []int `json:"bar"`
		Baz t03a  `json:"baz"`
	}

	src := t03{
		Foo: 1,
		Bar: []int{2, 3, 4},
		Baz: t03a{Aaa: 5},
	}

	dst := t03{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t03{
		Foo: 1,
		Bar: []int{2, 3, 4},
		Baz: t03a{Aaa: 5},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Typed -> Untyped
func Test4a(t *testing.T) {
	type t04a struct {
		Aaa int `json:"aaa"`
	}
	type t04 struct {
		Foo int   `json:"foo"`
		Bar []int `json:"bar"`
		Baz t04a  `json:"baz"`
	}

	src := t04{
		Foo: 1,
		Bar: []int{2, 3, 4},
		Baz: t04a{Aaa: 5},
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// Typed (fixed array) -> Untyped
func Test4b(t *testing.T) {
	type t04a struct {
		Aaa int `json:"aaa"`
	}
	type t04 struct {
		Foo int    `json:"foo"`
		Bar [2]int `json:"bar"`
		Baz t04a   `json:"baz"`
	}

	src := t04{
		Foo: 1,
		Bar: [2]int{2, 3},
		Baz: t04a{Aaa: 5},
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3},
		"baz": map[string]interface{}{
			"aaa": 5,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test5(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")

	src := map[string]interface{}{
		"foo": "SGVsbG8sIFdvcmxkIQ==",
		"bar": "1999-12-31T23:45:59.123456789Z",
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: []byte("Hello, World!"),
		Bar: wt1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test5b(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000Z", "1999-12-31T23:45:59.123Z")

	src := map[string]interface{}{
		"foo": "SGVsbG8sIFdvcmxkIQ==",
		"bar": "1999-12-31T23:45:59.123000Z",
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: []byte("Hello, World!"),
		Bar: wt1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test5c(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": "$$$",
		"bar": "1999-12-31T23:45:59.123456789Z",
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func Test5d(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": "SGVsbG8sIFdvcmxkIQ==",
		"bar": "A999-12-31T23:45:59.123456789Z",
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err == nil {
		t.Errorf("expect error\n")
	}
}

func Test6(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")

	src := t05{
		Foo: []byte("Hello, World!"),
		Bar: wt1,
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": "SGVsbG8sIFdvcmxkIQ==",
		"bar": "1999-12-31T23:45:59.123456789Z",
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test6b(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000Z", "1999-12-31T23:45:59.123Z")

	src := t05{
		Foo: []byte("Hello, World!"),
		Bar: wt1,
	}

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": "SGVsbG8sIFdvcmxkIQ==",
		"bar": "1999-12-31T23:45:59.123Z",
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test7(t *testing.T) {
	type t05 struct {
		Foo []byte    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	type t05r struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")

	src := t05{
		Foo: []byte("Hello, World!"),
		Bar: wt1,
	}

	dst := t05r{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05r{
		Foo: "SGVsbG8sIFdvcmxkIQ==",
		Bar: "1999-12-31T23:45:59.123456789Z",
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test8(t *testing.T) {
	type t05 struct {
		Foo []int     `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	type t05r struct {
		Foo string    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")

	src := t05{
		Foo: []int{1, 2, 3},
		Bar: wt1,
	}

	dst := t05r{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05r{
		Foo: "[1 2 3]",
		Bar: wt1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test9(t *testing.T) {
	type t05 struct {
		Foo int    `json:"foo"`
		Bar string `json:"bar"`
	}
	type t05r struct {
		Foo string    `json:"foo"`
		Bar time.Time `json:"bar"`
	}
	wt1, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "1999-12-31T23:45:59.123456789Z")

	src := t05{
		Foo: 123,
		Bar: "1999-12-31T23:45:59.123456789Z",
	}

	dst := t05r{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05r{
		Foo: "123",
		Bar: wt1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test10(t *testing.T) {
	type t05 struct {
		Foo float64 `json:"foo"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"nan": true,
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if !math.IsNaN(dst.Foo) {
		t.Errorf("dst: %v, want: %v\n", dst, math.NaN())
	}
}

func Test11(t *testing.T) {
	type t05 struct {
		Foo float64 `json:"foo"`
		Bar float64 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"inf": int(1),
		},
		"bar": map[string]interface{}{
			"inf": int(-1),
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: math.Inf(1),
		Bar: math.Inf(-1),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test11b(t *testing.T) {
	type t05 struct {
		Foo float64 `json:"foo"`
		Bar float64 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"inf": uint(1),
		},
		"bar": map[string]interface{}{
			"inf": uint(1),
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: math.Inf(1),
		Bar: math.Inf(1),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test11c(t *testing.T) {
	type t05 struct {
		Foo float64 `json:"foo"`
		Bar float64 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"inf": 1.0,
		},
		"bar": map[string]interface{}{
			"inf": -1.0,
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: math.Inf(1),
		Bar: math.Inf(-1),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test12(t *testing.T) {
	type t05 struct {
		Foo complex128 `json:"foo"`
		Bar complex128 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"re": float64(-12.34),
			"im": float64(56.78),
		},
		"bar": map[string]interface{}{
			"re": map[string]interface{}{"inf": -1},
			"im": map[string]interface{}{"inf": 1},
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: complex(-12.34, 56.78),
		Bar: complex(math.Inf(-1), math.Inf(1)),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test12b(t *testing.T) {
	type t05 struct {
		Foo complex128 `json:"foo"`
		Bar complex128 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"re": int(-12),
			"im": int(56),
		},
		"bar": map[string]interface{}{
			"re": map[string]interface{}{"inf": -1},
			"im": map[string]interface{}{"inf": 1},
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: complex(-12, 56),
		Bar: complex(math.Inf(-1), math.Inf(1)),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func Test12c(t *testing.T) {
	type t05 struct {
		Foo complex128 `json:"foo"`
		Bar complex128 `json:"bar"`
	}

	src := map[string]interface{}{
		"foo": map[string]interface{}{
			"re": uint(12),
			"im": uint(56),
		},
		"bar": map[string]interface{}{
			"re": map[string]interface{}{"inf": -1},
			"im": map[string]interface{}{"inf": 1},
		},
	}

	dst := t05{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t05{
		Foo: complex(12, 56),
		Bar: complex(math.Inf(-1), math.Inf(1)),
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}
