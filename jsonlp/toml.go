package jsonlp

import (
	"errors"
	"strconv"

	"github.com/shellyln/go-loose-json-parser/jsonlp/class"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

var (
	tomlParser ParserFn
)

func init() {
	tomlParser = tomlDocument()
}

func tomlTableKeyValuePair() ParserFn {
	return FlatGroup(
		objectKey(false),
		sp0NoLb(),
		erase(CharClass("=")),
		sp0NoLb(),
		First(
			FlatGroup(
				primitiveValue(),
				sp0NoLb(),
				First(
					erase(CharClass("\r\n", "\r", "\n")),
					LookAhead(End()),
				),
			),
			Indirect(listValue),
			Indirect(objectValue),
			Error("Expect object property value"),
		),
	)
}

func tomlArrayOfTable() ParserFn {
	return Trans(
		FlatGroup(
			erase(CharClass("[[")),
			First(
				FlatGroup(
					sp0NoLb(),
					objectKey(false),
					sp0NoLb(),
					erase(CharClass("]]")),
					sp0NoLb(),
				),
				Error("Expect array of table closing bracket ']]'"),
			),
			First(
				erase(CharClass("\r\n", "\r", "\n")),
				LookAhead(End()),
				Error("Expect line break or EOF"),
			),
			sp0(),
			Trans(
				ZeroOrMoreTimes(
					First(
						tomlTableKeyValuePair(),
					),
					sp0(),
				),
				tableTransformer,
				ChangeClassName(class.TomlArrayOfTable),
			),
		),
	)
}

func tomlTable() ParserFn {
	return Trans(
		FlatGroup(
			erase(CharClass("[")),
			First(
				FlatGroup(
					sp0NoLb(),
					objectKey(false),
					sp0NoLb(),
					erase(CharClass("]")),
					sp0NoLb(),
				),
				Error("Expect table closing bracket ']'"),
			),
			First(
				erase(CharClass("\r\n", "\r", "\n")),
				LookAhead(End()),
				Error("Expect line break or EOF"),
			),
			sp0(),
			Trans(
				ZeroOrMoreTimes(
					First(
						tomlTableKeyValuePair(),
					),
					sp0(),
				),
				tableTransformer,
				ChangeClassName(class.TomlTable),
			),
		),
	)
}

func tomlDocument() ParserFn {
	return Trans(
		FlatGroup(
			Start(),
			sp0(),
			OneOrMoreTimes(
				First(
					tomlTableKeyValuePair(),
					tomlArrayOfTable(),
					tomlTable(),
				),
				sp0(),
			),
			First(
				End(),
				Error("Expect terminatiion"),
			),
		),
		tableTransformer,
	)
}

// src: Loose TOML
//
// plafLb:
// Platform-dependent line break. (`Linebreak_Lf` | `Linebreak_CrLf` | `Linebreak_Cr`)
// Line break codes in multi-line string are replaced by this specified line break.
// (Excluding line breaks by escape sequences)
//
// interop:
// If Interop_JSON is set, replace NaN, Infinity, complex number by `{nan:true}`, `{inf:+/-1}`, `{re:re,im:im}`.
// If Interop_TOML is set, replace complex number by `{re:re,im:im}`.
// If Interop_JSON_AsNull is set, replace NaN, Infinity, complex number by null.
// If Interop_TOML_AsNull is set, replace complex number by null.
//
// parsed:
// nil | []any | map[string]any | float64 | int64 | uint64 | complex128 | string | bool | time.Time
func ParseTOML(s string, plafLb PlatformLinebreakType, interop InteropType) (interface{}, error) {
	ctx := *NewStringParserContext(s)
	opts := parseOptions{
		interop:           interop,
		platformLinebreak: "\n",
		isTOML:            true,
	}
	switch plafLb {
	case Linebreak_CrLf:
		opts.platformLinebreak = "\r\n"
	case Linebreak_Cr:
		opts.platformLinebreak = "\r"
	}
	ctx.Tag = opts

	out, err := tomlParser(ctx)
	if err != nil {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return nil, errors.New(
			err.Error() +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource)
	}

	if out.MatchStatus == MatchStatus_Matched {
		return out.AstStack[0].Value, nil
	} else {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return nil, errors.New(
			"Parse failed" +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource)
	}
}
