package jsonlp

import (
	"errors"
	"strconv"

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
		erase(sp0NoLb()),
		erase(CharClass("=")),
		erase(sp0NoLb()),
		First(
			FlatGroup(
				primitiveValue(),
				erase(sp0NoLb()),
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
					erase(sp0NoLb()),
					objectKey(false),
					erase(sp0NoLb()),
					erase(CharClass("]]")),
					erase(sp0NoLb()),
				),
				Error("Expect array of table closing bracket ']]'"),
			),
			First(
				erase(CharClass("\r\n", "\r", "\n")),
				LookAhead(End()),
				Error("Expect line break or EOF"),
			),
			erase(sp0()),
			Trans(
				ZeroOrMoreTimes(
					First(
						tomlTableKeyValuePair(),
					),
					erase(sp0()),
				),
				tableTransformer,
				ChangeClassName("TomlArrayOfTable"),
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
					erase(sp0NoLb()),
					objectKey(false),
					erase(sp0NoLb()),
					erase(CharClass("]")),
					erase(sp0NoLb()),
				),
				Error("Expect table closing bracket ']'"),
			),
			First(
				erase(CharClass("\r\n", "\r", "\n")),
				LookAhead(End()),
				Error("Expect line break or EOF"),
			),
			erase(sp0()),
			Trans(
				ZeroOrMoreTimes(
					First(
						tomlTableKeyValuePair(),
					),
					erase(sp0()),
				),
				tableTransformer,
				ChangeClassName("TomlTable"),
			),
		),
	)
}

func tomlDocument() ParserFn {
	return Trans(
		FlatGroup(
			Start(),
			erase(sp0()),
			OneOrMoreTimes(
				First(
					tomlTableKeyValuePair(),
					tomlArrayOfTable(),
					tomlTable(),
				),
				erase(sp0()),
			),
			First(
				End(),
				Error("Expect terminatiion"),
			),
		),
		tableTransformer,
	)
}

// src:     TOML
// interop: If true, replace NaN, Infinity by null
// parsed:  nil | []any | map[string]any | float64 | string | bool | time.Time
func ParseTOML(s string, interop bool) (interface{}, error) {
	ctx := *NewStringParserContext(s)
	ctx.Tag = parseOptions{interop: interop, isTOML: true}

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
