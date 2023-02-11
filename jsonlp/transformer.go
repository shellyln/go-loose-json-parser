package jsonlp

import (
	"strconv"

	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

func radixNumberTransformer(prefix string, radix int) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		switch asts[1].Value.(string) {
		case "p", "P":
			asts, err := Concat(ctx, asts)
			if err != nil {
				return nil, err
			}
			v, err := strconv.ParseFloat(prefix+asts[0].Value.(string), 64)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Float",
				Type:      AstType_Float,
				Value:     v,
			}}, nil
		case "s64", "S64":
			v, err := strconv.ParseUint(asts[0].Value.(string), radix, 64)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Int",
				Type:      AstType_Int,
				Value:     int64(v),
			}}, nil
		case "u64", "U64":
			v, err := strconv.ParseUint(asts[0].Value.(string), radix, 64)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Uint",
				Type:      AstType_Uint,
				Value:     v,
			}}, nil
		default:
			v, err := strconv.ParseInt(asts[0].Value.(string), radix, 64)
			if err != nil {
				return nil, err
			}
			return AstSlice{{
				ClassName: "Float",
				Type:      AstType_Float,
				Value:     float64(v),
			}}, nil
		}
	}
}

func decimalNumberTransformer(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	switch asts[1].Value.(string) {
	case "s64", "S64":
		v, err := strconv.ParseInt(asts[0].Value.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return AstSlice{{
			ClassName: "Int",
			Type:      AstType_Int,
			Value:     int64(v),
		}}, nil
	case "u64", "U64":
		v, err := strconv.ParseUint(asts[0].Value.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return AstSlice{{
			ClassName: "Uint",
			Type:      AstType_Uint,
			Value:     v,
		}}, nil
	default:
		v, err := strconv.ParseFloat(asts[0].Value.(string), 64)
		if err != nil {
			return nil, err
		}
		return AstSlice{{
			ClassName: "Float",
			Type:      AstType_Float,
			Value:     v,
		}}, nil
	}
}

func tableTransformer(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	length := len(asts)
	v := make(map[string]interface{})

	lastRefs := make(map[string]*map[string]interface{})
	lastRefs[""] = &v

	for i := 0; i < length; i += 2 {
		switch w := asts[i].Value.(type) {
		case string:
			// Simple identifier
			merged := false

			if asts[i+1].ClassName == "TomlArrayOfTable" {
				var a2 []map[string]interface{}
				if tmp, ok := v[w].([]map[string]interface{}); ok {
					a2 = tmp
				} else {
					// Register (case of v[w] == nil) or overwrite
					// NOTE: If overwrite, it is invalid TOML
					a2 = make([]map[string]interface{}, 0, 8)
					v[w] = a2
				}
				if m1, ok := asts[i+1].Value.(map[string]interface{}); ok {
					dottedKey := makeDottedKeyForSimpleName(w)
					lastRefs[dottedKey] = &m1
					a2 = append(a2, m1)
					v[w] = a2
				}
			} else {
				if m1, ok := asts[i+1].Value.(map[string]interface{}); ok {
					dottedKey := makeDottedKeyForSimpleName(w)
					lastRefs[dottedKey] = &m1
					if m2, ok := v[w].(map[string]interface{}); ok {
						// Merge redefined table
						// NOTE: Possibly an invalid TOML (except in cases such as `[x.y.z] ... [x]`)
						for xKey, xVal := range m1 {
							m2[xKey] = xVal
						}
						merged = true
					}
				}
				if !merged {
					v[w] = asts[i+1].Value
				}
			}

		case []string:
			// Dotted identifier
			for j, key := range w {
				if j == len(w)-1 {
					merged := false
					table := lastRefs[makeDottedKey(w, j)]

					if asts[i+1].ClassName == "TomlArrayOfTable" {
						var a2 []map[string]interface{}
						if tmp, ok := (*table)[key].([]map[string]interface{}); ok {
							a2 = tmp
						} else {
							// Register (case of (*table)[key] == nil) or overwrite
							// NOTE: If overwrite, it is invalid TOML
							a2 = make([]map[string]interface{}, 0, 8)
							(*table)[key] = a2
						}
						if m1, ok := asts[i+1].Value.(map[string]interface{}); ok {
							dottedKey := makeDottedKey(w, j+1)
							lastRefs[dottedKey] = &m1
							a2 = append(a2, m1)
							(*table)[key] = a2
						}
					} else {
						if m1, ok := asts[i+1].Value.(map[string]interface{}); ok {
							dottedKey := makeDottedKey(w, j+1)
							lastRefs[dottedKey] = &m1
							if m2, ok := (*table)[key].(map[string]interface{}); ok {
								// Merge redefined table
								// NOTE: Possibly an invalid TOML (except in cases such as `[x.y.z] ... [x.y]`)
								for xKey, xVal := range m1 {
									m2[xKey] = xVal
								}
								merged = true
							}
						}
						if !merged {
							(*table)[key] = asts[i+1].Value
						}
					}
				} else {
					prev := lastRefs[makeDottedKey(w, j)]
					dottedKey := makeDottedKey(w, j+1)
					if _, ok := lastRefs[dottedKey]; ok {
						// Already registerd to lastRefs
					} else {
						// Not registerd to lastRefs
						if cur, ok := (*prev)[key]; ok {
							switch next := cur.(type) {
							case map[string]interface{}:
								// Register
								lastRefs[dottedKey] = &next
							default:
								// Overwrite
								// NOTE: it is invalid TOML
								table := make(map[string]interface{})
								(*prev)[key] = table
								lastRefs[dottedKey] = &table
							}
						} else {
							// Append
							table := make(map[string]interface{})
							(*prev)[key] = table
							lastRefs[dottedKey] = &table
						}
					}
				}
			}
		}
	}
	return AstSlice{{
		ClassName: "Object",
		Type:      AstType_Any,
		Value:     v,
	}}, nil
}
