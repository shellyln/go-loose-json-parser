package jsonlp

import (
	"errors"
	"strconv"

	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

var (
	jsonParser ParserFn
)

func init() {
	jsonParser = jsonDocument()
}

func listValue() ParserFn {
	return Trans(
		FlatGroup(
			erase((Seq("["))),
			sp0(),
			ZeroOrOnce(
				FlatGroup(
					First(
						primitiveValue(),
						Indirect(listValue),
						Indirect(objectValue),
					),
					sp0(),
				),
				ZeroOrMoreTimes(
					erase((Seq(","))),
					sp0(),
					First(
						primitiveValue(),
						Indirect(listValue),
						Indirect(objectValue),
						LookAhead(Seq("]")),
						FlatGroup(
							sp0(),
							Error("Expect array closing parenthesis ')' or value"),
						),
					),
					sp0(),
				),
			),
			ZeroOrOnce(
				erase((Seq(","))),
				sp0(),
			),
			First(
				erase((Seq("]"))),
				FlatGroup(
					sp0(),
					Error("Expect array closing parenthesis ')'"),
				),
			),
			sp0(),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			length := len(asts)
			v := make([]interface{}, length)
			for i := 0; i < length; i++ {
				v[i] = asts[i].Value
			}
			return AstSlice{{
				ClassName: "Array",
				Type:      AstType_Any,
				Value:     v,
			}}, nil
		},
	)
}

func objectKey(allowLb bool) ParserFn {
	return Trans(
		First(
			dottedIdentifier(allowLb),
			stringValue(),
			identifier(),
		),
	)
}

func objectKeyValuePair() ParserFn {
	return FlatGroup(
		objectKey(true),
		sp0(),
		erase(First(CharClass(":"), CharClass("=>"), CharClass("="))),
		sp0(),
		First(
			primitiveValue(),
			Indirect(listValue),
			Indirect(objectValue),
			Error("Expect object property value"),
		),
		sp0(),
	)
}

func objectValue() ParserFn {
	return Trans(
		FlatGroup(
			erase((Seq("{"))),
			sp0(),
			ZeroOrOnce(
				objectKeyValuePair(),
				ZeroOrMoreTimes(
					erase((Seq(","))),
					sp0(),
					First(
						objectKeyValuePair(),
						LookAhead(Seq("}")),
						FlatGroup(
							sp0(),
							Error("Expect object closing bracket '}' or key-value pair"),
						),
					),
				),
				ZeroOrOnce(
					erase((Seq(","))),
					sp0(),
				),
			),
			First(
				erase((Seq("}"))),
				FlatGroup(
					sp0(),
					Error("Expect object closing bracket '}'"),
				),
			),
			sp0(),
		),
		tableTransformer,
	)
}

func jsonDocument() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			primitiveValue(),
			listValue(),
			objectValue(),
		),
		sp0(),
		First(
			End(),
			Error("Expect terminatiion"),
		),
	)
}

// src:     Loose JSON
// interop: If true, replace NaN, Infinity by null
// parsed:  nil | []any | map[string]any | float64 | int64 | uint64 | complex128 | string | bool | time.Time
func Parse(s string, interop bool) (interface{}, error) {
	ctx := *NewStringParserContext(s)
	ctx.Tag = parseOptions{interop: interop}

	out, err := jsonParser(ctx)
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
