package jsonlp

import (
	"errors"
	"strconv"
	"time"

	. "github.com/shellyln/takenoco/base"
	"github.com/shellyln/takenoco/extra"
	. "github.com/shellyln/takenoco/string"
)

var (
	documentParser ParserFn
)

func init() {
	documentParser = document()
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Whitespaces
func sp0() ParserFn {
	return erase(ZeroOrMoreTimes(First(Whitespace(), comment())))
}

func lineComment() ParserFn {
	return erase(FlatGroup(
		Seq("//"),
		FlatGroup(
			ZeroOrMoreTimes(CharClassN("\r\n", "\n", "\r")),
			First(
				CharClass("\r\n", "\n", "\r"),
				LookAhead(End()),
			),
		),
	))
}

func hashLineComment() ParserFn {
	return erase(FlatGroup(
		Seq("#"),
		FlatGroup(
			ZeroOrMoreTimes(CharClassN("\r\n", "\n", "\r")),
			First(
				CharClass("\r\n", "\n", "\r"),
				LookAhead(End()),
			),
		),
	))
}

func blockComment() ParserFn {
	return erase(FlatGroup(
		Seq("/*"),
		ZeroOrMoreTimes(CharClassN("*/")),
		First(
			Seq("*/"),
			Error("An unexpected termination has appeared in the block comment."),
		),
	))
}

func comment() ParserFn {
	return First(lineComment(), hashLineComment(), blockComment())
}

func trueValue() ParserFn {
	return FlatGroup(
		erase(Seq("true")),
		WordBoundary(),
		Zero(trueAst),
	)
}

func falseValue() ParserFn {
	return FlatGroup(
		erase(Seq("false")),
		WordBoundary(),
		Zero(falseAst),
	)
}

func boolValue() ParserFn {
	return First(
		trueValue(),
		falseValue(),
	)
}

func nullValue() ParserFn {
	return FlatGroup(
		erase(CharClass("null", "undefined")),
		WordBoundary(),
		Zero(nilAst),
	)
}

func positiveInfinityValue() ParserFn {
	return FlatGroup(
		erase(FlatGroup(
			ZeroOrOnce(Seq("+")),
			Seq("Infinity"),
		)),
		WordBoundary(),
		Zero(positiveInfinityAst),
	)
}

func negativeInfinityValue() ParserFn {
	return FlatGroup(
		erase(Seq("-Infinity")),
		WordBoundary(),
		Zero(negativeInfinityAst),
	)
}

func nanValue() ParserFn {
	return FlatGroup(
		erase(Seq("NaN")),
		WordBoundary(),
		Zero(nanAst),
	)
}

func numberValue() ParserFn {
	return First(
		FlatGroup(
			First(
				Trans(
					FlatGroup(erase(SeqI("0b")), extra.BinaryNumberStr()),
					ParseIntRadix(2),
					castIntToFloat,
				),
				Trans(
					FlatGroup(erase(SeqI("0o")), extra.OctalNumberStr()),
					ParseIntRadix(8),
					castIntToFloat,
				),
				Trans(
					FlatGroup(erase(SeqI("0x")), extra.HexNumberStr()),
					ParseIntRadix(16),
					castIntToFloat,
				),
				Trans(
					extra.FloatNumberStr(),
					ParseFloat,
					ChangeClassName("Float"),
				),
				Trans(
					extra.IntegerNumberStr(),
					ParseInt,
					castIntToFloat,
				),
			),
			WordBoundary(),
		),
		positiveInfinityValue(),
		negativeInfinityValue(),
		nanValue(),
	)
}

func stringLiteralInner(cc string, multiline bool) ParserFn {
	return FlatGroup(
		erase(Seq(cc)),
		ZeroOrMoreTimes(
			First(
				FlatGroup(
					erase(Seq("\\")),
					First(
						CharClass("\\", "'", "\"", "`"),
						replaceStr(CharClass("n", "N"), "\n"),
						replaceStr(CharClass("r", "R"), "\r"),
						replaceStr(CharClass("v", "V"), "\v"),
						replaceStr(CharClass("t", "T"), "\t"),
						replaceStr(CharClass("b", "B"), "\b"),
						replaceStr(CharClass("f", "F"), "\f"),
						Trans(
							FlatGroup(
								erase(CharClass("u")),
								Repeat(Times{Min: 4, Max: 4}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								erase(CharClass("u{")),
								Repeat(Times{Min: 1, Max: 6}, HexNumber()),
								erase(CharClass("}")),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								erase(CharClass("x")),
								Repeat(Times{Min: 2, Max: 2}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								Repeat(Times{Min: 3, Max: 3}, OctNumber()),
							),
							ParseIntRadix(8),
							StringFromInt,
						),
					),
				),
				If(multiline,
					OneOrMoreTimes(CharClassN(cc, "\\")),
					OneOrMoreTimes(
						First(
							FlatGroup(
								CharClass("\r", "\n"),
								Error("An unexpected newline has appeared in the string literal."),
							),
							CharClassN(cc, "\\"),
						),
					),
				),
			),
		),
		First(
			FlatGroup(End(), Error("An unexpected termination has appeared in the string literal.")),
			erase(Seq(cc)),
		),
	)
}

func stringValue() ParserFn {
	return Trans(
		First(
			stringLiteralInner("\"", false),
			stringLiteralInner("'", false),
			stringLiteralInner("`", true),
		),
		Concat,
		ChangeClassName("String"),
	)
}

func dateValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.DateStr(),
			WordBoundary(),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			value := asts[len(asts)-1].Value.(string)
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Date",
				Type:      AstType_Any,
				Value:     t.UTC(),
			}}, nil
		},
	)
}

func dateTimeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.DateTimeStr(),
			WordBoundary(),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			// TODO: BUG: Cannot parse years with negative values or years greater than or equal to 10000.
			value := asts[len(asts)-1].Value.(string)
			t, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", value)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "DateTime",
				Type:      AstType_Any,
				Value:     t.UTC(),
			}}, nil
		},
	)
}

func timeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.TimeStr(),
			WordBoundary(),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			value := "1970-01-01T" + asts[len(asts)-1].Value.(string) + "+00:00"
			t, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", value)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Time",
				Type:      AstType_Any,
				Value:     t.UTC(),
			}}, nil
		},
	)
}

func primitiveValue() ParserFn {
	return First(
		stringValue(),
		boolValue(),
		nullValue(),
		timeValue(),
		dateTimeValue(),
		dateValue(),
		numberValue(),
	)
}

func listValue() ParserFn {
	return Trans(
		FlatGroup(
			erase((Seq("["))),
			erase(sp0()),
			ZeroOrOnce(
				FlatGroup(
					First(
						primitiveValue(),
						Indirect(listValue),
						Indirect(objectValue),
					),
					erase(sp0()),
				),
				ZeroOrMoreTimes(
					erase((Seq(","))),
					erase(sp0()),
					First(
						primitiveValue(),
						Indirect(listValue),
						Indirect(objectValue),
						LookAhead(Seq("]")),
						FlatGroup(
							erase(sp0()),
							Error("Expect array closing parenthesis ')' or value"),
						),
					),
					erase(sp0()),
				),
			),
			ZeroOrOnce(
				erase((Seq(","))),
				erase(sp0()),
			),
			First(
				erase((Seq("]"))),
				FlatGroup(
					erase(sp0()),
					Error("Expect array closing parenthesis ')'"),
				),
			),
			erase(sp0()),
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

func identifier() ParserFn {
	return Trans(
		extra.UnicodeIdentifierStr(),
		ChangeClassName("Idenitifier"),
	)
}

func objectKey() ParserFn {
	return Trans(
		First(
			stringValue(),
			identifier(),
		),
	)
}

func objectKeyValuePair() ParserFn {
	return FlatGroup(
		objectKey(),
		sp0(),
		erase(First(CharClass(":"), CharClass("=>"))),
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
			erase(sp0()),
			ZeroOrOnce(
				objectKeyValuePair(),
				ZeroOrMoreTimes(
					erase((Seq(","))),
					erase(sp0()),
					First(
						objectKeyValuePair(),
						LookAhead(Seq("}")),
						FlatGroup(
							erase(sp0()),
							Error("Expect object closing bracket '}' or key-value pair"),
						),
					),
				),
				ZeroOrOnce(
					erase((Seq(","))),
					erase(sp0()),
				),
			),
			First(
				erase((Seq("}"))),
				FlatGroup(
					erase(sp0()),
					Error("Expect object closing bracket '}'"),
				),
			),
			erase(sp0()),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			length := len(asts)
			v := make(map[string]interface{})
			for i := 0; i < length; i += 2 {
				v[asts[i].Value.(string)] = asts[i+1].Value
			}
			return AstSlice{{
				ClassName: "Object",
				Type:      AstType_Any,
				Value:     v,
			}}, nil
		},
	)
}

func document() ParserFn {
	return FlatGroup(
		Start(),
		erase(sp0()),
		First(
			primitiveValue(),
			listValue(),
			objectValue(),
			erase(sp0()),
		),
		End(),
	)
}

func Parse(s string) (interface{}, error) {
	out, err := documentParser(*NewStringParserContext(s))
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