package jsonlp

import (
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

func tomlSingleLineLiteralString() ParserFn {
	return FlatGroup(
		erase(Seq("'")),
		ZeroOrMoreTimes(
			First(
				OneOrMoreTimes(
					First(
						FlatGroup(
							CharClass("\r", "\n"),
							Error("An unexpected newline has appeared in the string literal."),
						),
						CharClassN("'"),
					),
				),
			),
		),
		First(
			FlatGroup(End(), Error("An unexpected termination has appeared in the string literal.")),
			erase(Seq("'")),
		),
	)
}

func tomlMultiLineLiteralString() ParserFn {
	return FlatGroup(
		erase(Seq("'''")),
		ZeroOrOnce(erase(CharClass("\r\n", "\r", "\n"))),
		ZeroOrMoreTimes(
			First(
				FlatGroup(
					CharClass("'"),
					LookAhead(CharClass("'''")),
				),
				OneOrMoreTimes(CharClassN("'''")),
			),
		),
		First(
			FlatGroup(End(), Error("An unexpected termination has appeared in the string literal.")),
			erase(Seq("'''")),
		),
	)
}

func tomlSingleLineBasicString() ParserFn {
	return FlatGroup(
		erase(Seq("\"")),
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
								Repeat(Times{Min: 6, Max: 6}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
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
				OneOrMoreTimes(
					First(
						FlatGroup(
							CharClass("\r", "\n"),
							Error("An unexpected newline has appeared in the string literal."),
						),
						CharClassN("\"", "\\"),
					),
				),
			),
		),
		First(
			FlatGroup(End(), Error("An unexpected termination has appeared in the string literal.")),
			erase(Seq("\"")),
		),
	)
}

func tomlMultiLineBasicString() ParserFn {
	return FlatGroup(
		erase(Seq("\"\"\"")),
		ZeroOrOnce(erase(CharClass("\r\n", "\r", "\n"))),
		ZeroOrMoreTimes(
			First(
				FlatGroup(
					erase(Seq("\\")),
					First(
						erase(FlatGroup(
							CharClass("\r\n", "\r", "\n"),
							sp0(),
						)),
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
								Repeat(Times{Min: 6, Max: 6}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
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
				OneOrMoreTimes(
					First(
						FlatGroup(
							CharClass("\""),
							LookAhead(CharClass("\"\"\"")),
						),
						CharClassN("\"\"\"", "\\"),
					),
				),
			),
		),
		First(
			FlatGroup(End(), Error("An unexpected termination has appeared in the string literal.")),
			erase(Seq("\"\"\"")),
		),
	)
}

func tomlStringValueInner() ParserFn {
	return First(
		tomlMultiLineBasicString(),
		tomlSingleLineBasicString(),
		tomlMultiLineLiteralString(),
		tomlSingleLineLiteralString(),
		stringLiteralInner("`", true),
	)
}

func tomlStringValue() ParserFn {
	return Trans(
		FlatGroup(
			tomlStringValueInner(),
			ZeroOrMoreTimes(
				sp1NoLb(),
				tomlStringValueInner(),
			),
		),
		Concat,
		ChangeClassName("String"),
	)
}
