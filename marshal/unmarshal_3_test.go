package marshal_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/shellyln/go-loose-json-parser/marshal"
)

// recursive (map -> map)
func TestRecursive1(t *testing.T) {
	src := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4},
		"baz": nil,
	}
	src["bar"] = append(src["bar"].([]interface{}), src)
	src["baz"] = src

	dst := map[string]interface{}{}
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := map[string]interface{}{
		"foo": 1,
		"bar": []interface{}{2, 3, 4, nil},
		"baz": nil,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

// recursive (struct -> struct)
func TestRecursive2(t *testing.T) {
	type t02 struct {
		Foo int  `json:"foo"`
		Bar *t02 `json:"bar"`
		Baz *t02 `json:"baz"`
	}

	src := t02{
		Foo: 1,
		Bar: &t02{
			Foo: 2,
		},
	}
	src.Bar.Bar = &src
	src.Baz = src.Bar
	src.Bar.Baz = src.Baz

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, &marshal.MarshalOptions{
		NoCopyUnexportedFields: true,
	}); err != nil {
		t.Errorf("%v\n", err)
	}
	dstByte, _ := json.Marshal(dst)
	dstStr := string(dstByte)

	// NOTE: `src` is passed by-val
	want := t02{
		Foo: 1,
		Bar: &t02{
			Foo: 2,
			Bar: &t02{
				Foo: 1,
			},
		},
		Baz: &t02{
			Foo: 2,
			Bar: &t02{
				Foo: 1,
			},
		},
	}
	wantByte, _ := json.Marshal(want)
	wantStr := string(wantByte)

	if src.Foo != dst.Foo {
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if dst.Bar == nil || dst.Baz == nil {
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if src.Bar == dst.Bar {
		// should be copied
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if src.Bar.Bar == dst.Bar.Bar {
		// should be copied
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if src.Baz == dst.Baz {
		// should be copied
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if src.Baz.Bar == dst.Baz.Bar {
		// should be copied
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
	if dstStr != wantStr {
		t.Errorf("dst: %v, want: %v\n", dstStr, wantStr)
	}
}

// nil (struct -> struct)
func TestRecursive3(t *testing.T) {
	type t02 struct {
		Foo int  `json:"foo"`
		Bar *t02 `json:"bar"`
		Baz *t02 `json:"baz"`
	}

	src := t02{
		Foo: 1,
	}

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, &marshal.MarshalOptions{
		NoCopyUnexportedFields: true,
	}); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t02{
		Foo: 1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
}

// nil (struct -> struct)
func TestRecursive4(t *testing.T) {
	type t02 struct {
		Foo int  `json:"foo"`
		Bar *t02 `json:"bar"`
		Baz *t02 `json:"baz"`
	}

	src := t02{
		Foo: 1,
		Bar: &t02{
			Foo: 2,
		},
	}

	dst := t02{}
	if err := marshal.Unmarshal(src, &dst, &marshal.MarshalOptions{
		NoCopyUnexportedFields: true,
	}); err != nil {
		t.Errorf("%v\n", err)
	}

	want := t02{
		Foo: 1,
		Bar: &t02{
			Foo: 2,
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("src: %v, dst: %v, want: %v\n", src, dst, want)
	}
}

type mytypeTestCustom1 int

func (s mytypeTestCustom1) MarshalLp() (interface{}, error) {
	return s * 2, nil
}

func TestCustom1(t *testing.T) {
	src := mytypeTestCustom1(1234)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(2468)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestCustom1b(t *testing.T) {
	src := mytypeTestCustom1(1234)

	var dst *int
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst == nil {
		t.Errorf("dst: %v\n", dst)
	}
	if *dst != int(2468) {
		t.Errorf("dst: %v\n", *dst)
	}
}

func TestCustom1c(t *testing.T) {
	srcTmp := mytypeTestCustom1(1234)
	src := &srcTmp

	var dst int
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst != int(2468) {
		t.Errorf("dst: %v\n", dst)
	}
}

type mytypeTestCustom2 int

func (p *mytypeTestCustom2) MarshalLp() (interface{}, error) {
	if p != nil {
		return *p * 2, nil
	} else {
		return 0, nil
	}
}

func TestCustom2(t *testing.T) {
	src := mytypeTestCustom2(1234)

	dst := int(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := int(2468)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestCustom2b(t *testing.T) {
	src := mytypeTestCustom2(1234)

	var dst *int
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst == nil {
		t.Errorf("dst: %v\n", dst)
	}
	if *dst != int(2468) {
		t.Errorf("dst: %v\n", *dst)
	}
}

func TestCustom2c(t *testing.T) {
	srcTmp := mytypeTestCustom2(1234)
	src := &srcTmp

	var dst int
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst != int(2468) {
		t.Errorf("dst: %v\n", dst)
	}
}

type mytypeTestCustom3 int

func (p *mytypeTestCustom3) UnmarshalLp(from interface{}) error {
	switch v := from.(type) {
	case int:
		*p = mytypeTestCustom3(v) * 2
	}
	return nil
}

func TestCustom3(t *testing.T) {
	src := int(1234)

	dst := mytypeTestCustom3(0)
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	want := mytypeTestCustom3(2468)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestCustom3b(t *testing.T) {
	src := int(1234)

	dstTmp := mytypeTestCustom3(0)
	dst := &dstTmp
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst == nil {
		t.Errorf("dst: %v\n", dst)
	}
	if *dst != mytypeTestCustom3(2468) {
		t.Errorf("dst: %v\n", *dst)
	}
}

func TestCustom3c(t *testing.T) {
	srcTmp := int(1234)
	src := &srcTmp

	var dst mytypeTestCustom3
	if err := marshal.Unmarshal(src, &dst, nil); err != nil {
		t.Errorf("%v\n", err)
	}

	if dst != mytypeTestCustom3(2468) {
		t.Errorf("dst: %v\n", dst)
	}
}
