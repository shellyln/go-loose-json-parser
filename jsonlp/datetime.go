package jsonlp

import (
	"github.com/shellyln/go-loose-json-parser/jsonlp/class"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

// Customized `extra. DateTimeStr()`
//
// Parse the ISO 8601 datetime string.
// (yyyy-MM-ddThh:mmZ , ... , yyyy-MM-ddThh:mm:ss.fffffffffZ)
// (yyyy-MM-ddThh:mm+00:00 , ... , yyyy-MM-ddThh:mm:ss.fffffffff+00:00)
func dateTimeStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(Seq("-")),
			Repeat(Times{Min: 4, Max: -1}, Number()),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '1'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '3'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			First(
				Seq("T"),
				// TOML allows Datetime format with date and time delimited by space (RFC 3339 section 5.6)
				FlatGroup(
					erase(Seq(" ")),
					Zero(Ast{
						Type:  AstType_String,
						Value: "T",
					}),
				),
			),
			CharRange(RuneRange{Start: '0', End: '2'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq(":"),
			CharRange(RuneRange{Start: '0', End: '5'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			First(
				FlatGroup(
					Seq(":"),
					CharRange(RuneRange{Start: '0', End: '6'}),
					CharRange(RuneRange{Start: '0', End: '9'}),
					First(
						FlatGroup(
							Seq("."),
							Trans(
								Repeat(Times{Min: 1, Max: 9}, // 3: milli, 6: micro, 9: nano
									CharRange(RuneRange{Start: '0', End: '9'}),
								),
								Concat,
								func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
									return AstSlice{{
										Type:  AstType_String,
										Value: (asts[len(asts)-1].Value.(string) + "000000000")[0:9],
									}}, nil
								},
							),
						),
						Zero(Ast{
							Type:  AstType_String,
							Value: ".000000000",
						}),
					),
				),
				Zero(Ast{
					Type:  AstType_String,
					Value: ":00.000000000",
				}),
			),
			First(
				FlatGroup(
					erase(Seq("Z")),
					Zero(Ast{
						Type:  AstType_String,
						Value: "+00:00",
					}),
				),
				FlatGroup(
					CharClass("+", "-"),
					Repeat(Times{Min: 2, Max: 2},
						CharRange(RuneRange{Start: '0', End: '9'}),
					),
					Seq(":"),
					CharRange(RuneRange{Start: '0', End: '5'}),
					CharRange(RuneRange{Start: '0', End: '9'}),
				),
				FlatGroup(
					// TOML allows Datetime format without timezone
					Zero(Ast{
						Type:  AstType_String,
						Value: "+00:00",
					}),
				),
			),
		),
		Concat,
		ChangeClassName(class.DateTimeStr),
	)
}
