package marshal_test

import (
	"reflect"
	"testing"

	"github.com/shellyln/go-loose-json-parser/jsonlp"
	"github.com/shellyln/go-loose-json-parser/marshal"
)

func TestTag1(t *testing.T) {
	type config struct {
		Addr string `json:"addr"`
	}
	type response struct {
		Config config `json:"config"`
	}

	parsed, err := jsonlp.Parse(`{
        // comment
        config: {
            addr: '127.0.0.1',
        }
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst response
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := response{
		Config: config{
			Addr: "127.0.0.1",
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestTag2(t *testing.T) {
	type config struct {
		Addr string
	}
	type response struct {
		Config config
	}

	parsed, err := jsonlp.Parse(`{
        // comment
        config: {
            addr: '127.0.0.1',
        }
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst response
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := response{
		Config: config{
			Addr: "",
		},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull1(t *testing.T) {
	type s1 struct {
		Foo int  `json:"foo"`
		Bar *int `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: null,
        bar: null,
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: 0,
		Bar: nil,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull2(t *testing.T) {
	type s1 struct {
		Foo int `json:"foo"`
		Bar int `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: "",
        bar: "1",
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: 0,
		Bar: 1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull3(t *testing.T) {
	type s1 struct {
		Foo uint `json:"foo"`
		Bar uint `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: "",
        bar: "1",
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: 0,
		Bar: 1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull4(t *testing.T) {
	type s1 struct {
		Foo float64 `json:"foo"`
		Bar float64 `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: "",
        bar: "1",
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: 0,
		Bar: 1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull5(t *testing.T) {
	type s1 struct {
		Foo complex128 `json:"foo"`
		Bar complex128 `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: "",
        bar: "1",
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: 0,
		Bar: 1,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}

func TestZeroOrNull6(t *testing.T) {
	type s1 struct {
		Foo bool `json:"foo"`
		Bar bool `json:"bar"`
	}

	parsed, err := jsonlp.Parse(`{
        foo: "",
        bar: "1",
    }`, jsonlp.Interop_None)

	if err != nil {
		t.Errorf("Parse: error = %v\n", err)
		return
	}

	var dst s1
	if err := marshal.Unmarshal(parsed, &dst, nil); err != nil {
		t.Errorf("Unmarshal: error = %v\n", err)
		return
	}

	want := s1{
		Foo: false,
		Bar: true,
	}
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst: %v, want: %v\n", dst, want)
	}
}
