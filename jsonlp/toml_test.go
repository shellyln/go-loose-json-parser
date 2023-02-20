package jsonlp_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shellyln/go-loose-json-parser/jsonlp"
)

func runMatrixTomlParse(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbg {
				fmt.Println("")
			}
			got, err := jsonlp.ParseTOML(tt.args.s, tt.args.plafLb, tt.args.interop)
			if tt.dbg {
				fmt.Println("")
			}
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

func TestTomlParse1(t *testing.T) {
	tests := []testMatrixItem{{
		name: "t1-0a",
		args: args{s: `
		[x.y.z]
		b=2
		[x]
		a=1
		`},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"a": float64(1),
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"b": float64(2),
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t1-0b",
		args: args{s: `
		[x]
		a=1
		[x.y.z]
		b=2
		`},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"a": float64(1),
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"b": float64(2),
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t1-1a",
		args: args{s: `
		foo=1
		bar.baz=2
		`},
		want: map[string]interface{}{
			"foo": float64(1),
			"bar": map[string]interface{}{
				"baz": float64(2),
			},
		},
		wantErr: false,
	}, {
		name: "t1-1b",
		args: args{s: `
		foo = 1
		bar . baz = 2
		`},
		want: map[string]interface{}{
			"foo": float64(1),
			"bar": map[string]interface{}{
				"baz": float64(2),
			},
		},
		wantErr: false,
	}, {
		name: "t1-1c",
		args: args{s: `
		"foo" = 1
		"bar" . "baz" = 2
		`},
		want: map[string]interface{}{
			"foo": float64(1),
			"bar": map[string]interface{}{
				"baz": float64(2),
			},
		},
		wantErr: false,
	}, {
		name: "t1-2a",
		args: args{s: `
		fff=999

		[x]
		zz=0

		[x.y.z]
		aa=1
		bb=2

		[x]
		cc=3
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"zz": float64(0),
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"aa": float64(1),
						"bb": float64(2),
					},
				},
				"cc": float64(3),
			},
		},
		wantErr: false,
	}, {
		name: "t1-2a2",
		args: args{s: `
		"fff" = -999.0

		[ "x" ]
		"zz" = 0.0

		[ "x" . y . "z" ]
		"aa" = -1.0
		"bb" = 2.0

		[ x ]
		# merge it
		cc = 3.0
		# overwrite previous [x.y.z]
		"y" . z . "dd" = 4
		`},
		want: map[string]interface{}{
			"fff": float64(-999),
			"x": map[string]interface{}{
				"zz": float64(0),
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"dd": float64(4),
					},
				},
				"cc": float64(3),
			},
		},
		wantErr: false,
	}, {
		name: "t1-2b",
		args: args{s: `
		fff=999

		[x]
		zz=0

		[x.y.z]
		aa=1
		bb=2
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"zz": float64(0),
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"aa": float64(1),
						"bb": float64(2),
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t1-2c",
		args: args{s: `
		fff=999

		[x.y.z]
		aa=1
		bb=2

		[x]
		cc=3
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"z": map[string]interface{}{
						"aa": float64(1),
						"bb": float64(2),
					},
				},
				"cc": float64(3),
			},
		},
		wantErr: false,
	}, {
		name: "t1-3a",
		args: args{s: `
		fff=999

		[x.y]
		zz=0

		[x.y.z]
		aa=1
		bb=2

		[x.y]
		cc=3
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"zz": float64(0),
					"z": map[string]interface{}{
						"aa": float64(1),
						"bb": float64(2),
					},
					"cc": float64(3),
				},
			},
		},
		wantErr: false,
	}, {
		// https://toml.io/en/v1.0.0#keys
		name: "t1-4a",
		args: args{s: `
		apple.type = "fruit"
		apple.skin = "thin"
		apple.color = "red"

		orange.type = "fruit"
		orange.skin = "thick"
		orange.color = "orange"
		`},
		want: map[string]interface{}{
			"apple": map[string]interface{}{
				"type":  "fruit",
				"skin":  "thin",
				"color": "red",
			},
			"orange": map[string]interface{}{
				"type":  "fruit",
				"skin":  "thick",
				"color": "orange",
			},
		},
		wantErr: false,
		dbg:     true,
	}, {
		// https://toml.io/en/v1.0.0#keys
		name: "t1-4b",
		args: args{s: `
		apple.type = "fruit"
		orange.type = "fruit"
		
		apple.skin = "thin"
		orange.skin = "thick"
		
		apple.color = "red"
		orange.color = "orange"
		`},
		want: map[string]interface{}{
			"apple": map[string]interface{}{
				"type":  "fruit",
				"skin":  "thin",
				"color": "red",
			},
			"orange": map[string]interface{}{
				"type":  "fruit",
				"skin":  "thick",
				"color": "orange",
			},
		},
		wantErr: false,
		dbg:     true,
	}, {
		// https://toml.io/en/v1.0.0#keys
		name: "t1-5a",
		args: args{s: `
		key = "value1"
		bare_key = "value2"
		bare-key = "value3"
		1234 = "value4"
		-3.14 = "value5"
		`},
		want: map[string]interface{}{
			"key":      "value1",
			"bare_key": "value2",
			"bare-key": "value3",
			"1234":     "value4",
			"-3": map[string]interface{}{
				"14": "value5",
			},
		},
		wantErr: false,
		dbg:     true,
	}, {
		name: "t1-5a",
		args: args{s: `
		foo = {
			key: "value1",
			bare_key: "value2",
			bare-key: "value3",
			1234: "value4",
			-3.14: "value5",
		}
		`},
		want: map[string]interface{}{
			"foo": map[string]interface{}{
				"key":      "value1",
				"bare_key": "value2",
				"bare-key": "value3",
				"1234":     "value4",
				"-3": map[string]interface{}{
					"14": "value5",
				},
			},
		},
		wantErr: false,
		dbg:     true,
	}}

	runMatrixTomlParse(t, tests)
}

func TestTomlParse2(t *testing.T) {
	tests := []testMatrixItem{{
		name: "t2-1a",
		args: args{s: `
		fff=999

		[[x]]
		aa=1
		bb=2
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": []map[string]interface{}{
				{
					"aa": float64(1),
					"bb": float64(2),
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-1b",
		args: args{s: `
		fff=999

		[[x]]
		aa=1
		bb=2
		[[x]]
		aa=3
		bb=4
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": []map[string]interface{}{
				{
					"aa": float64(1),
					"bb": float64(2),
				},
				{
					"aa": float64(3),
					"bb": float64(4),
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-2a",
		args: args{s: `
		fff=999

		[[x]]
		aa=1
		bb=2
		[x.y]
		cc=11
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": []map[string]interface{}{
				{
					"aa": float64(1),
					"bb": float64(2),
					"y": map[string]interface{}{
						"cc": float64(11),
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-2b",
		args: args{s: `
		fff=999

		[[x]]
		aa=1
		bb=2
		[x.y]
		cc=11
		[[x]]
		aa=3
		bb=4
		[x.y]
		cc=21
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": []map[string]interface{}{
				{
					"aa": float64(1),
					"bb": float64(2),
					"y": map[string]interface{}{
						"cc": float64(11),
					},
				},
				{
					"aa": float64(3),
					"bb": float64(4),
					"y": map[string]interface{}{
						"cc": float64(21),
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-3a",
		args: args{s: `
		fff=999

		[[x.y.z]]
		aa=1
		bb=2
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"z": []map[string]interface{}{
						{
							"aa": float64(1),
							"bb": float64(2),
						},
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-3b",
		args: args{s: `
		fff=999

		[[x.y.z]]
		aa=1
		bb=2
		[[x.y.z]]
		aa=3
		bb=4
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"z": []map[string]interface{}{
						{
							"aa": float64(1),
							"bb": float64(2),
						},
						{
							"aa": float64(3),
							"bb": float64(4),
						},
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-4a",
		args: args{s: `
		fff=999

		[[x.y.z]]
		aa=1
		bb=2
		[x.y.z.p]
		cc=11
		dd=12
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"z": []map[string]interface{}{
						{
							"aa": float64(1),
							"bb": float64(2),
							"p": map[string]interface{}{
								"cc": float64(11),
								"dd": float64(12),
							},
						},
					},
				},
			},
		},
		wantErr: false,
	}, {
		name: "t2-4b",
		args: args{s: `
		fff=999

		[[x.y.z]]
		aa=1
		bb=2
		[x.y.z.p]
		cc=11
		dd=12
		[[x.y.z]]
		aa=3
		bb=4
		[x.y.z.p]
		cc=21
		dd=22
		`},
		want: map[string]interface{}{
			"fff": float64(999),
			"x": map[string]interface{}{
				"y": map[string]interface{}{
					"z": []map[string]interface{}{
						{
							"aa": float64(1),
							"bb": float64(2),
							"p": map[string]interface{}{
								"cc": float64(11),
								"dd": float64(12),
							},
						},
						{
							"aa": float64(3),
							"bb": float64(4),
							"p": map[string]interface{}{
								"cc": float64(21),
								"dd": float64(22),
							},
						},
					},
				},
			},
		},
		wantErr: false,
	}}

	runMatrixTomlParse(t, tests)
}

func TestTomlParse3(t *testing.T) {
	tests := []testMatrixItem{{
		name: "t3-1a",
		args: args{s: `
		str1 = "foo" "bar"
		str2 = "foo"  'bar'
		`},
		want: map[string]interface{}{
			"str1": "foobar",
			"str2": "foobar",
		},
		wantErr: false,
	}, {
		// https://toml.io/en/v1.0.0#string
		name: "t3-2a",
		args: args{s: `
		str1 = "The quick brown fox jumps over the lazy dog."

		str2 = """
The quick brown \


		fox jumps over \
		  the lazy dog."""

		str3 = """\
			   The quick brown \
			   fox jumps over \
			   the lazy dog.\
			   """
		`},
		want: map[string]interface{}{
			"str1": "The quick brown fox jumps over the lazy dog.",
			"str2": "The quick brown fox jumps over the lazy dog.",
			"str3": "The quick brown fox jumps over the lazy dog.",
		},
		wantErr: false,
	}, {
		// https://toml.io/en/v1.0.0#string
		name: "t3-3a",
		args: args{s: `
		quot15 = '''Here are fifteen quotation marks: """""""""""""""'''
		apos15 = "Here are fifteen apostrophes: '''''''''''''''"
		str1 = ''''That,' she said, 'is still pointless.'''
		str2 = ''''That,' she said, 'is still pointless.''''
		str3 = """"That," she said, "is still pointless."""
		str4 = """"That," she said, "is still pointless.""""
		`},
		want: map[string]interface{}{
			"quot15": `Here are fifteen quotation marks: """""""""""""""`,
			"apos15": `Here are fifteen apostrophes: '''''''''''''''`,
			"str1":   "'That,' she said, 'is still pointless.",
			"str2":   "'That,' she said, 'is still pointless.'",
			"str3":   `"That," she said, "is still pointless.`,
			"str4":   `"That," she said, "is still pointless."`,
		},
		wantErr: false,
	}, {
		name: "t3-4a",
		args: args{s: `
		str1 = '''
The quick\nbrown
fox jumps over
the lazy dog.'''
		`},
		want: map[string]interface{}{
			"str1": "The quick\\nbrown\nfox jumps over\nthe lazy dog.",
		},
		wantErr: false,
	}, {
		name: "t3-5a",
		args: args{s: `
		str1 = """
The quick\nbrown
fox jumps over
the lazy dog."""
		`},
		want: map[string]interface{}{
			"str1": "The quick\nbrown\nfox jumps over\nthe lazy dog.",
		},
		wantErr: false,
	}, {
		name: "t3-6a",
		args: args{s: `
		str1 = '''
The quick\nbrown
fox jumps over
the lazy dog.'''
		`, plafLb: jsonlp.Linebreak_CrLf},
		want: map[string]interface{}{
			"str1": "The quick\\nbrown\r\nfox jumps over\r\nthe lazy dog.",
		},
		wantErr: false,
	}, {
		name: "t3-7a",
		args: args{s: `
		str1 = """
The quick\nbrown
fox jumps over
the lazy dog."""
		`, plafLb: jsonlp.Linebreak_CrLf},
		want: map[string]interface{}{
			"str1": "The quick\nbrown\r\nfox jumps over\r\nthe lazy dog.",
		},
		wantErr: false,
	}, {
		name: "t3-8a",
		args: args{s: `
		str1 = '''
The quick\nbrown
fox jumps over
the lazy dog.'''
		`, plafLb: jsonlp.Linebreak_Cr},
		want: map[string]interface{}{
			"str1": "The quick\\nbrown\rfox jumps over\rthe lazy dog.",
		},
		wantErr: false,
	}, {
		name: "t3-9a",
		args: args{s: `
		str1 = """
The quick\nbrown
fox jumps over
the lazy dog."""
		`, plafLb: jsonlp.Linebreak_Cr},
		want: map[string]interface{}{
			"str1": "The quick\nbrown\rfox jumps over\rthe lazy dog.",
		},
		wantErr: false,
	}}

	runMatrixTomlParse(t, tests)
}

func TestTomlParse4(t *testing.T) {
	// Test Strict IEEE-754 +0/-0
	got, err := jsonlp.ParseTOML(`x = 0.0`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("0.0: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "0p-1074" {
		t.Errorf("0.0 not equals 0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = +0.0`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("+0.0: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "0p-1074" {
		t.Errorf("+0.0 not equals 0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = -0.0`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("-0.0: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "-0p-1074" {
		t.Errorf("-0.0 not equals -0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = +0`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("+0: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "0p-1074" {
		t.Errorf("+0 not equals 0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = -0`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("-0: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "-0p-1074" {
		t.Errorf("-0 not equals -0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = 0x0p-1074`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("0x0p-1074: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "0p-1074" {
		t.Errorf("0x0p-1074 not equals 0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = +0x0p-1074`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("+0x0p-1074: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "0p-1074" {
		t.Errorf("+0x0p-1074 not equals 0p-1074")
		return
	}

	got, err = jsonlp.ParseTOML(`x = -0x0p-1074`, jsonlp.Linebreak_Lf, jsonlp.Interop_None)
	if err != nil {
		t.Errorf("-0x0p-1074: Parse() error = %v", err)
		return
	}
	if fmt.Sprintf("%b", got.(map[string]interface{})["x"]) != "-0p-1074" {
		t.Errorf("-0x0p-1074 not equals -0p-1074")
		return
	}
}
