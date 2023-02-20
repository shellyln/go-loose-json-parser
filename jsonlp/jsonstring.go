package jsonlp

import (
	"github.com/shellyln/go-loose-json-parser/jsonlp/class"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

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
						replaceStr(CharClass("\r\n", "\r", "\n"), ""),
						Any(),
					),
				),
				If(multiline,
					OneOrMoreTimes(
						First(
							Trans(
								CharClass("\r\n", "\r", "\n"),
								func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
									return AstSlice{Ast{
										Type:  AstType_String,
										Value: ctx.Tag.(parseOptions).platformLinebreak,
									}}, nil
								},
							),
							CharClassN(cc, "\\"),
						),
					),
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

func jsonStringValue() ParserFn {
	return Trans(
		First(
			tomlMultiLineBasicString(),
			tomlMultiLineLiteralString(),
			stringLiteralInner("\"", false),
			stringLiteralInner("'", false),
			stringLiteralInner("`", true),
		),
		Concat,
		ChangeClassName(class.String),
	)
}
