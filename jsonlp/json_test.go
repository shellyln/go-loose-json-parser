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
	dbg     bool
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

func TestJsonParse1(t *testing.T) {
	tests := []testMatrixItem{{
		name:    "j1-1a",
		args:    args{s: `null`},
		want:    nil,
		wantErr: false,
	}, {
		name:    "j1-1b",
		args:    args{s: ` null `},
		want:    nil,
		wantErr: false,
	}, {
		name:    "j1-2",
		args:    args{s: `undefined`},
		want:    nil,
		wantErr: false,
	}, {
		name:    "j1-3",
		args:    args{s: `true`},
		want:    true,
		wantErr: false,
	}, {
		name:    "j1-4",
		args:    args{s: `false`},
		want:    false,
		wantErr: false,
	}, {
		name:    "j1-5a",
		args:    args{s: `0`},
		want:    float64(0),
		wantErr: false,
	}, {
		name:    "j1-5b",
		args:    args{s: `-12.34`},
		want:    float64(-12.34),
		wantErr: false,
	}, {
		name:    "j1-5c1",
		args:    args{s: `1234`},
		want:    float64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c2",
		args:    args{s: `+1234`},
		want:    float64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c3",
		args:    args{s: `-1234`},
		want:    float64(-1234),
		wantErr: false,
	}, {
		name:    "j1-5c4",
		args:    args{s: `1234s64`},
		want:    int64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c5",
		args:    args{s: `+1234s64`},
		want:    int64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c6",
		args:    args{s: `-1234s64`},
		want:    int64(-1234),
		wantErr: false,
	}, {
		name:    "j1-5c7",
		args:    args{s: `1234u64`},
		want:    uint64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c8",
		args:    args{s: `+1234u64`},
		want:    uint64(1234),
		wantErr: false,
	}, {
		name:    "j1-5c9",
		args:    args{s: `-0s64`},
		want:    int64(0),
		wantErr: false,
	}, {
		name:    "j1-5c10",
		args:    args{s: `9223372036854775807s64`},
		want:    int64(9223372036854775807),
		wantErr: false,
	}, {
		name:    "j1-5c10b",
		args:    args{s: `9_223_372_036_854_775_807s64`},
		want:    int64(9223372036854775807),
		wantErr: false,
	}, {
		name:    "j1-5c11",
		args:    args{s: `-9223372036854775808s64`},
		want:    int64(-9223372036854775808),
		wantErr: false,
	}, {
		name:    "j1-5c11b",
		args:    args{s: `-9_223_372_036_854_775_808_s64`},
		want:    int64(-9223372036854775808),
		wantErr: false,
	}, {
		name:    "j1-5c12",
		args:    args{s: `-1s64`},
		want:    int64(-1),
		wantErr: false,
	}, {
		name:    "j1-5c13",
		args:    args{s: `9223372036854775807u64`},
		want:    uint64(9223372036854775807),
		wantErr: false,
	}, {
		name:    "j1-5c14",
		args:    args{s: `9223372036854775808u64`},
		want:    uint64(9223372036854775808),
		wantErr: false,
	}, {
		name:    "j1-5c15",
		args:    args{s: `18446744073709551615u64`},
		want:    uint64(18446744073709551615),
		wantErr: false,
	}, {
		name:    "j1-5c15b",
		args:    args{s: `18_446_744_073_709_551_615_u64`},
		want:    uint64(18446744073709551615),
		wantErr: false,
	}, {
		name:    "j1-5c16a",
		args:    args{s: `9007199254740991`},
		want:    float64(9007199254740991),
		wantErr: false,
	}, {
		name:    "j1-5c16b",
		args:    args{s: `9_007_199_254_740_991`},
		want:    float64(9007199254740991),
		wantErr: false,
	}, {
		name:    "j1-5c16c",
		args:    args{s: `-9007199254740991`},
		want:    float64(-9007199254740991),
		wantErr: false,
	}, {
		name:    "j1-5c16d",
		args:    args{s: `-9_007_199_254_740_991`},
		want:    float64(-9007199254740991),
		wantErr: false,
	}, {
		name:    "j1-5c17",
		args:    args{s: `-123_456.789_012`},
		want:    float64(-123456.789012),
		wantErr: false,
	}, {
		name:    "j1-5d1",
		args:    args{s: `0b1100`},
		want:    float64(12),
		wantErr: false,
	}, {
		name:    "j1-5d2",
		args:    args{s: `0b0011`},
		want:    float64(3),
		wantErr: false,
	}, {
		name:    "j1-5e",
		args:    args{s: `0o0040`},
		want:    float64(32),
		wantErr: false,
	}, {
		name:    "j1-5f1a",
		args:    args{s: `0x0080`},
		want:    float64(128),
		wantErr: false,
	}, {
		name:    "j1-5f1b",
		args:    args{s: `0x0080s64`},
		want:    int64(128),
		wantErr: false,
	}, {
		name:    "j1-5f1c",
		args:    args{s: `0x0080u64`},
		want:    uint64(128),
		wantErr: false,
	}, {
		name:    "j1-5f1d",
		args:    args{s: `+0x0080`},
		wantErr: true,
	}, {
		name:    "j1-5f1e",
		args:    args{s: `-0x0080`},
		wantErr: true,
	}, {
		name:    "j1-5f1f",
		args:    args{s: `+0x0080s64`},
		wantErr: true,
	}, {
		name:    "j1-5f1g",
		args:    args{s: `-0x0080s64`},
		wantErr: true,
	}, {
		name:    "j1-5f1h",
		args:    args{s: `+0x0080u64`},
		wantErr: true,
	}, {
		name:    "j1-5f1i",
		args:    args{s: `-0x0080u64`},
		wantErr: true,
	}, {
		name:    "j1-5f2",
		args:    args{s: `0x7fffffffffffffffs64`},
		want:    int64(9223372036854775807),
		wantErr: false,
	}, {
		name:    "j1-5f3",
		args:    args{s: `0x8000000000000000s64`},
		want:    int64(-9223372036854775808),
		wantErr: false,
	}, {
		name:    "j1-5f4",
		args:    args{s: `0xffffffffffffffffs64`},
		want:    int64(-1),
		wantErr: false,
	}, {
		name:    "j1-5f4b",
		args:    args{s: `0x_ffff_ffff_ffff_ffff_s64`},
		want:    int64(-1),
		wantErr: false,
	}, {
		name:    "j1-5f5",
		args:    args{s: `0x7fffffffffffffffu64`},
		want:    uint64(9223372036854775807),
		wantErr: false,
	}, {
		name:    "j1-5f6",
		args:    args{s: `0x8000000000000000u64`},
		want:    uint64(9223372036854775808),
		wantErr: false,
	}, {
		name:    "j1-5f7",
		args:    args{s: `0xffffffffffffffffu64`},
		want:    uint64(18446744073709551615),
		wantErr: false,
	}, {
		name:    "j1-5f7b",
		args:    args{s: `0x_ffff_ffff_ffff_ffff_u64`},
		want:    uint64(18446744073709551615),
		wantErr: false,
	}, {
		name:    "j1-5f8",
		args:    args{s: `0x1p-2`},
		want:    float64(0.25),
		wantErr: false,
	}, {
		name:    "j1-5f9a",
		args:    args{s: `0x1.Fp+0`},
		want:    float64(1.9375),
		wantErr: false,
	}, {
		name:    "j1-5f9b",
		args:    args{s: `0x1.F000_0000_p+0`},
		want:    float64(1.9375),
		wantErr: false,
	}, {
		name:    "j1-5f9c",
		args:    args{s: `0x1._F000_0000_p+0`},
		want:    float64(1.9375),
		wantErr: false,
	}, {
		name:    "j1-5f9d",
		args:    args{s: `0x1._F000_0000_p0`},
		want:    float64(1.9375),
		wantErr: false,
	}, {
		name:    "j1-5f9e",
		args:    args{s: `0x1._F000_0000_p-0`},
		want:    float64(1.9375),
		wantErr: false,
	}, {
		name:    "j1-5f9f",
		args:    args{s: `0x1.8p+1`},
		want:    float64(3),
		wantErr: false,
	}, {
		name:    "j1-5f9g",
		args:    args{s: `0x1.8p+0`},
		want:    float64(1.5),
		wantErr: false,
	}, {
		name:    "j1-5f9h",
		args:    args{s: `0x1.8p0`},
		want:    float64(1.5),
		wantErr: false,
	}, {
		name:    "j1-5f9i",
		args:    args{s: `0x1.8p-1`},
		want:    float64(0.75),
		wantErr: false,
	}, {
		name:    "j1-5f10a",
		args:    args{s: `0X_1FFFP-16`},
		want:    float64(0.1249847412109375),
		wantErr: false,
	}, {
		name:    "j1-5f10b",
		args:    args{s: `0X_1F_FF_P-1_6_`},
		want:    float64(0.1249847412109375),
		wantErr: false,
	}, {
		name:    "j1-5f11",
		args:    args{s: `0X1FFFP-16`},
		want:    float64(0.1249847412109375),
		wantErr: false,
	}, {
		name:    "j1-5f12",
		args:    args{s: `-0X1FFFP-16`},
		want:    float64(-0.1249847412109375),
		wantErr: false,
	}, {
		name:    "j1-5g1a",
		args:    args{s: `-9.5e-3`},
		want:    float64(-0.0095),
		wantErr: false,
	}, {
		name:    "j1-5g1b",
		args:    args{s: `-9.5e3`},
		want:    float64(-9500),
		wantErr: false,
	}, {
		name:    "j1-5g1c",
		args:    args{s: `-9.5e+3`},
		want:    float64(-9500),
		wantErr: false,
	}, {
		name:    "j1-6a",
		args:    args{s: `Infinity`},
		want:    math.Inf(1),
		wantErr: false,
	}, {
		name:    "j1-6b",
		args:    args{s: `+Infinity`},
		want:    math.Inf(1),
		wantErr: false,
	}, {
		name:    "j1-6c",
		args:    args{s: `-Infinity`},
		want:    math.Inf(-1),
		wantErr: false,
	}, {
		name:    "j1-7a",
		args:    args{s: `"abc"`},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "j1-7b",
		args:    args{s: `'abc'`},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "j1-7c",
		args:    args{s: "`abc`"},
		want:    "abc",
		wantErr: false,
	}, {
		name:    "j1-7d",
		args:    args{s: `"abc\ndef"`},
		want:    "abc\ndef",
		wantErr: false,
	}, {
		name:    "j1-7e",
		args:    args{s: `"\u0048\u0065\u006c\u006c\u006f\u002c\u0020\u0057\u006f\u0072\u006c\u0064\u0021"`},
		want:    "Hello, World!",
		wantErr: false,
	}, {
		name:    "j1-7f",
		args:    args{s: `"\x48\x65\x6c\x6c\x6f\x2c\x20\x57\x6f\x72\x6c\x64\x21"`},
		want:    "Hello, World!",
		wantErr: false,
	}, {
		name:    "j1-8a",
		args:    args{s: `{}`},
		want:    map[string]interface{}{},
		wantErr: false,
	}, {
		name:    "j1-8b",
		args:    args{s: `{"abc":"def"}`},
		want:    map[string]interface{}{"abc": "def"},
		wantErr: false,
	}, {
		name:    "j1-8c",
		args:    args{s: `{"abc":"def",}`},
		want:    map[string]interface{}{"abc": "def"},
		wantErr: false,
	}, {
		name:    "j1-8d",
		args:    args{s: `{"abc":"def","ghi":123}`},
		want:    map[string]interface{}{"abc": "def", "ghi": float64(123)},
		wantErr: false,
	}, {
		name:    "j1-8e",
		args:    args{s: `{"abc":"def","ghi":123,}`},
		want:    map[string]interface{}{"abc": "def", "ghi": float64(123)},
		wantErr: false,
	}, {
		name:    "j1-9a",
		args:    args{s: `[]`},
		want:    []interface{}{},
		wantErr: false,
	}, {
		name:    "j1-9b",
		args:    args{s: `["abc"]`},
		want:    []interface{}{"abc"},
		wantErr: false,
	}, {
		name:    "j1-9c",
		args:    args{s: `["abc",]`},
		want:    []interface{}{"abc"},
		wantErr: false,
	}, {
		name:    "j1-9d",
		args:    args{s: `["abc","def"]`},
		want:    []interface{}{"abc", "def"},
		wantErr: false,
	}, {
		name:    "j1-9e",
		args:    args{s: `["abc","def",]`},
		want:    []interface{}{"abc", "def"},
		wantErr: false,
	}, {
		name: "j1-10a",
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
		name: "j1-10b",
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
		name: "j1-11a",
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
		name: "j1-11b",
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
		name: "j1-12a",
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
		name: "j1-13a",
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
		name: "j1-13b",
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
		name: "j1-13c",
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
	}, {
		name:    "j1-14a",
		args:    args{s: `1.23-34.5i`},
		want:    complex(1.23, -34.5),
		wantErr: false,
	}, {
		name:    "j1-14b",
		args:    args{s: ` 1.23 - 34.5i `},
		want:    complex(1.23, -34.5),
		wantErr: false,
	}, {
		name:    "j1-14b2",
		args:    args{s: ` 1.23 - 34.5_i `},
		want:    complex(1.23, -34.5),
		wantErr: false,
	}, {
		name:    "j1-14c",
		args:    args{s: `-1.23+34.5i`},
		want:    complex(-1.23, +34.5),
		wantErr: false,
	}, {
		name:    "j1-14d",
		args:    args{s: ` -1.23 + 34.5i `},
		want:    complex(-1.23, +34.5),
		wantErr: false,
	}, {
		name:    "j1-14e",
		args:    args{s: `123s64-345u64i`},
		want:    complex(123, -345),
		wantErr: false,
	}, {
		name:    "j1-14f",
		args:    args{s: `123u64-345s64i`},
		want:    complex(123, -345),
		wantErr: false,
	}, {
		name:    "j1-14g",
		args:    args{s: `0xffs64-0x7fu64i`},
		want:    complex(255, -127),
		wantErr: false,
	}, {
		name:    "j1-14h",
		args:    args{s: `0xffu64-0x7fs64i`},
		want:    complex(255, -127),
		wantErr: false,
	}, {
		name:    "j1-14i",
		args:    args{s: `1.23e+1-34.5e-1i`},
		want:    complex(12.3, -3.45),
		wantErr: false,
	}, {
		name:    "j1-14j",
		args:    args{s: `0x1.8p+1-0x1.8p-1i`},
		want:    complex(3, -0.75),
		wantErr: false,
		// }, {
		// 	name:    "j1-14k",
		// 	args:    args{s: `NaN-NaNi`},
		// 	want:    complex(math.NaN(), math.NaN()),
		// 	wantErr: false,
	}, {
		name:    "j1-14l1",
		args:    args{s: `Infinity-Infinityi`},
		want:    complex(math.Inf(1), math.Inf(-1)),
		wantErr: false,
	}, {
		name:    "j1-14l2",
		args:    args{s: `Infinity - Infinity_i`},
		want:    complex(math.Inf(1), math.Inf(-1)),
		wantErr: false,
	}, {
		name:    "j1-14m1",
		args:    args{s: `Infinity--Infinityi`},
		want:    complex(math.Inf(1), math.Inf(1)),
		wantErr: false,
	}, {
		name:    "j1-14m2",
		args:    args{s: `Infinity - -Infinity_i`},
		want:    complex(math.Inf(1), math.Inf(1)),
		wantErr: false,
	}}

	runMatrixParse(t, tests)
}
