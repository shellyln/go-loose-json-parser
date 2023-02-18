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
