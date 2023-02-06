package jsonlp_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/shellyln/go-loose-json-parser/jsonlp"
)

type args struct {
	s string
}

type testMatrixItem struct {
	name    string
	args    args
	want    interface{}
	wantErr bool
}

func runMatrixParse(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonlp.Parse(tt.args.s, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("%v: Parse() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v: Parse() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []testMatrixItem{{
		name:    "1",
		args:    args{s: `null`},
		want:    nil,
		wantErr: false,
	}, {
		name:    "2",
		args:    args{s: `undefined`},
		want:    nil,
		wantErr: false,
	}, {
		name:    "3",
		args:    args{s: `true`},
		want:    true,
		wantErr: false,
	}, {
		name:    "4",
		args:    args{s: `false`},
		want:    false,
		wantErr: false,
	}, {
		name:    "5a",
		args:    args{s: `0`},
		want:    float64(0),
		wantErr: false,
	}, {
		name:    "5b",
		args:    args{s: `-12.34`},
		want:    float64(-12.34),
		wantErr: false,
	}, {
		name:    "5c",
		args:    args{s: `1234`},
		want:    float64(1234),
		wantErr: false,
	}, {
		name:    "5d",
		args:    args{s: `0b1100`},
		want:    float64(12),
		wantErr: false,
	}, {
		name:    "5e",
		args:    args{s: `0o0040`},
		want:    float64(32),
		wantErr: false,
	}, {
		name:    "5f",
		args:    args{s: `0x0080`},
		want:    float64(128),
		wantErr: false,
	}, {
		name:    "5g",
		args:    args{s: `-9.5e-3`},
		want:    float64(-0.0095),
		wantErr: false,
	}, {
		name:    "6a",
		args:    args{s: `Infinity`},
		want:    math.Inf(1),
		wantErr: false,
	}, {
		name:    "6b",
		args:    args{s: `+Infinity`},
		want:    math.Inf(1),
		wantErr: false,
	}, {
		name:    "6c",
		args:    args{s: `-Infinity`},
		want:    math.Inf(-1),
		wantErr: false,
	}, {
		name:    "7a",
		args:    args{s: `"abc"`},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "7b",
		args:    args{s: `'abc'`},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "7c",
		args:    args{s: "`abc`"},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "7d",
		args:    args{s: `"abc\ndef"`},
		want:    "abc\ndef",
		wantErr: false,
	}, {
		name:    "7e",
		args:    args{s: `"\u0048\u0065\u006c\u006c\u006f\u002c\u0020\u0057\u006f\u0072\u006c\u0064\u0021"`},
		want:    "Hello, World!",
		wantErr: false,
	}, {
		name:    "7f",
		args:    args{s: `"\x48\x65\x6c\x6c\x6f\x2c\x20\x57\x6f\x72\x6c\x64\x21"`},
		want:    "Hello, World!",
		wantErr: false,
	}, {
		name:    "8a",
		args:    args{s: `{}`},
		want:    map[string]interface{}{},
		wantErr: false,
	}, {
		name:    "8b",
		args:    args{s: `{"abc":"def"}`},
		want:    map[string]interface{}{"abc": "def"},
		wantErr: false,
	}, {
		name:    "8c",
		args:    args{s: `{"abc":"def",}`},
		want:    map[string]interface{}{"abc": "def"},
		wantErr: false,
	}, {
		name:    "8d",
		args:    args{s: `{"abc":"def","ghi":123}`},
		want:    map[string]interface{}{"abc": "def", "ghi": float64(123)},
		wantErr: false,
	}, {
		name:    "8e",
		args:    args{s: `{"abc":"def","ghi":123,}`},
		want:    map[string]interface{}{"abc": "def", "ghi": float64(123)},
		wantErr: false,
	}, {
		name:    "9a",
		args:    args{s: `[]`},
		want:    []interface{}{},
		wantErr: false,
	}, {
		name:    "9b",
		args:    args{s: `["abc"]`},
		want:    []interface{}{"abc"},
		wantErr: false,
	}, {
		name:    "9c",
		args:    args{s: `["abc",]`},
		want:    []interface{}{"abc"},
		wantErr: false,
	}, {
		name:    "9d",
		args:    args{s: `["abc","def"]`},
		want:    []interface{}{"abc", "def"},
		wantErr: false,
	}, {
		name:    "9e",
		args:    args{s: `["abc","def",]`},
		want:    []interface{}{"abc", "def"},
		wantErr: false,
	}, {
		name: "10a",
		args: args{s: `{亜a_$:[1,/**/2,'3abd',4,-Infinity,null,0x12_34,undefined,true,false,` +
			`2020-01-02,18:20:30.001,2020-01-02T18:20:30.001Z,{_c1:1,$c1:-1,` +
			`'d'=>2,'dd':3,"ddd"=4,` + "`dddd`" + ` : 5,f1:["eee"],},],bb亜:-12.34,}`},
		want: map[string]interface{}{
			"亜a_$": []interface{}{
				float64(1), float64(2), "3abd", float64(4), math.Inf(-1), nil, float64(4660), nil,
				true, false,
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(1970, time.January, 1, 18, 20, 30, 1000000, time.UTC),
				time.Date(2020, time.January, 2, 18, 20, 30, 1000000, time.UTC),
				map[string]interface{}{
					"_c1":  float64(1),
					"$c1":  float64(-1),
					"d":    float64(2),
					"dd":   float64(3),
					"ddd":  float64(4),
					"dddd": float64(5),
					"f1":   []interface{}{"eee"},
				},
			},
			"bb亜": float64(-12.34),
		},
		wantErr: false,
	}, {
		name: "10b",
		args: args{s: ` /**/ {
			亜a_$ : [
				// line comment
				1 , /* block comment */ 2 , '3abd' , 4 , -Infinity , null , 0x12_34 , undefined ,
				true , false ,
				2020-01-02 ,
				# line comment
				18:20:30.001 ,
				2020-01-02T18:20:30.001Z ,
				{ /**/
					_c1 : 1 /**/ ,
					$c1 /**/ : -1  ,
					'd' => /**/ 2 ,
					'dd' : 3 , // comment
					"ddd" = 4 ,
					` + "`dddd`" + ` : 5 ,
					f1 : [ /**/  "eee" /**/  ] , /**/
				} , /**/ 
			] , bb亜 : -12.34 ,
		} `},
		want: map[string]interface{}{
			"亜a_$": []interface{}{
				float64(1), float64(2), "3abd", float64(4), math.Inf(-1), nil, float64(4660), nil,
				true, false,
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(1970, time.January, 1, 18, 20, 30, 1000000, time.UTC),
				time.Date(2020, time.January, 2, 18, 20, 30, 1000000, time.UTC),
				map[string]interface{}{
					"_c1":  float64(1),
					"$c1":  float64(-1),
					"d":    float64(2),
					"dd":   float64(3),
					"ddd":  float64(4),
					"dddd": float64(5),
					"f1":   []interface{}{"eee"},
				},
			},
			"bb亜": float64(-12.34),
		},
		wantErr: false,
	}, {
		name: "11a",
		args: args{s: ` /**/ {
			亜a_$ : [
				// line comment
				1 , /* block comment */ 2 , '3abd' , 4 aa , -Infinity , null , 0x12_34 , undefined ,
				true , false ,
				2020-01-02 ,
				# line comment
				18:20:30.001 ,
				2020-01-02T18:20:30.001Z ,
				{ /**/
					_c1 : 1 /**/ ,
					$c1 /**/ : -1  ,
					'd' => /**/ 2 ,
					'dd' : 3 , // comment
					"ddd" = 4 ,
					` + "`dddd`" + ` : 5 ,
					f1 : [ /**/  "eee" /**/  ] , /**/
				} , /**/ 
			] , bb亜 : -12.34 ,
		} `},
		want:    nil,
		wantErr: true,
	}, {
		name: "11b",
		args: args{s: ` /**/ {
			亜a_$ : [
				// line comment
				1 , /* block comment */ 2 , '3abd' , 4 , -Infinity , null , 0x12_34 , undefined ,
				true , false ,
				2020-01-02 ,
				# line comment
				18:20:30.001 ,
				2020-01-02T18:20:30.001Z ,
				{ /**/
					_c1 : 1 /**/ ,
					$c1 /**/ : -1  ,
					'd' => /**/ 2 ,
					'dd' : 3 , // comment
					"ddd" = 4 aaa ,
					` + "`dddd`" + ` : 5 ,
					f1 : [ /**/  "eee" /**/  ] , /**/
				} , /**/ 
			] , bb亜 : -12.34 ,
		} `},
		want:    nil,
		wantErr: true,
	}, {
		name: "12a",
		args: args{s: `[
			123,          // -> float64(123)
			-123.45,      // -> float64(-123.45)
			-1.2345e+6,   // -> float64(-1234500)
			-123_456_789, // -> float64(-123456789)
			0x12345678,   // -> float64(305419896)
			0x1234_5678,  // -> float64(305419896)
			0o7654_3210,  // -> float64(16434824)
			0b0101_0101,  // -> float64(85)
			0x99999999,
			9007199254740991,
			-9007199254740991,
			0x1F_FFFF_FFFF_FFFF,
			9007199254740992,
			9007199254740993,
        ]`},
		want: []interface{}{
			float64(123),
			float64(-123.45),
			float64(-1234500),
			float64(-123456789),
			float64(305419896),
			float64(305419896),
			float64(16434824),
			float64(85),
			float64(2576980377),
			float64(9007199254740991), // MAX_SAFE_INTEGER
			float64(-9007199254740991),
			float64(9007199254740991),
			float64(9007199254740992), // MAX_SAFE_INTEGER+1
			float64(9007199254740992), // MAX_SAFE_INTEGER+1 == MAX_SAFE_INTEGER+2
		},
		wantErr: false,
	}, {
		name: "13a",
		args: args{s: `{foo.bar.baz=123}`},
		want: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": float64(123),
				},
			},
		},
		wantErr: false,
	}, {
		name: "13b",
		args: args{s: `{"foo"."bar"."baz"=123}`},
		want: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": float64(123),
				},
			},
		},
		wantErr: false,
	}, {
		name: "13c",
		args: args{s: `{foo.bar.baz=123, foo.bar.qux:234}`},
		want: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": float64(123),
					"qux": float64(234),
				},
			},
		},
		wantErr: false,
	}}

	runMatrixParse(t, tests)
}
