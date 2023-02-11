package jsonlp

import (
	"encoding/base64"
	"strings"

	. "github.com/shellyln/takenoco/base"
)

func setStr(s string) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		return AstSlice{{
			ClassName: "setStr",
			Type:      AstType_String,
			Value:     s,
		}}, nil
	}
}

func replaceStr(fn ParserFn, s string) ParserFn {
	return Trans(fn, setStr(s))
}

func makeDottedKey(name []string, length int) string {
	var keySb strings.Builder
	for i := 0; i < length; i++ {
		if i != 0 {
			keySb.WriteRune('.')
		}
		x := base64.StdEncoding.EncodeToString([]byte(name[i]))
		keySb.WriteString(x)
	}
	return keySb.String()
}

func makeDottedKeyForSimpleName(name string) string {
	var keySb strings.Builder
	x := base64.StdEncoding.EncodeToString([]byte(name))
	keySb.WriteString(x)
	return keySb.String()
}
