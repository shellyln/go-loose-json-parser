# Go Loose JSON + TOML parsers
Super loose JSON + TOML parsers for Go.

[![Test](https://github.com/shellyln/go-loose-json-parser/actions/workflows/test.yml/badge.svg)](https://github.com/shellyln/go-loose-json-parser/actions/workflows/test.yml)
[![release](https://img.shields.io/github/v/release/shellyln/go-loose-json-parser)](https://github.com/shellyln/go-loose-json-parser/releases)
[![Go version](https://img.shields.io/github/go-mod/go-version/shellyln/go-loose-json-parser)](https://github.com/shellyln/go-loose-json-parser)


## 🚀 Usage

### JSON
```go
package main

import (
    "fmt"
    "github.com/shellyln/go-loose-json-parser/jsonlp"
)

func main() {
    // src:     Loose JSON
    // interop: If true, replace NaN, Infinity by null
    // parsed:  nil | []any | map[string]any | float64 | string | bool | time.Time
    parsed, err := jsonlp.Parse(`{
        // comment
        config: {
            addr: '127.0.0.1',
        }
    }`, false)

    if err != nil {
        fmt.Printf("Parse: error = %v\n", err)
        return
    }

    fmt.Printf("Parsed = %v\n", parsed)
}
```

### TOML

```go
package main

import (
    "fmt"
    "github.com/shellyln/go-loose-json-parser/jsonlp"
)

func main() {
    // src:     Loose TOML
    // interop: If true, replace NaN, Infinity by null
    // parsed:  nil | []any | map[string]any | float64 | string | bool | time.Time
    parsed, err := jsonlp.ParseTOML(`
    # comment
    [config]
    addr = '127.0.0.1'
    `, false)

    if err != nil {
        fmt.Printf("Parse: error = %v\n", err)
        return
    }

    fmt.Printf("Parsed = %v\n", parsed)
}
```

## 🪄 Example

### JSON

```js
# Hash comment
{
    // Line comment
    /* Block comment */

    // Object keys can be enclosed in either double-quote, single-quote, back-quote.
    // It is also allowed not to enclose them.

    "foo": [
        123,          // -> float64(123)
        -123.45,      // -> float64(-123.45)
        -1.2345e+6,   // -> float64(-1234500)
        -123_456_789, // -> float64(-123456789)
        0x12345678,   // -> float64(305419896)
        0x1234_5678,  // -> float64(305419896)
        0o7654_3210,  // -> float64(16434824)
        0b0101_0101,  // -> float64(85)
    ],

    'bar': null,      // -> nil
    baz: undefined,   // -> nil
    "qux": -Infinity, // -> -Inf // Infinity, +Infinity are also available
    "quux": NaN,      // -> NaN

    // Non-ASCII identifiers can also be used.
    あいうえお: {
        key: "value1",
        bare_key: "value2",
        bare-key: "value3", // Keys containing hyphens are allowed.
        1234: "value4",     // Keys starting with a number are allowed.
        -3.14: "value5",    // "-3": { "14": "value5" }
    },

    "corge": [
        // Escape sequences are available
        'Hello, World!\r\n',

        // Unicode escape sequences (\uXXXX and \u{XXXXXX}) are available
        "\u0048\u0065\u006c\u006c\u006f\u002c\u0020\u0057\u006f\u0072\u006c\u0064\u0021",

        // Byte escape sequence is also available
        "\x48\x65\x6c\x6c\x6f\x2c\x20\x57\x6f\x72\x6c\x64\x21",

        // Multiline string literal
        `Tempor adipiscing amet velit ipsum diam ut ea takimata lorem gubergren sed laoreet.
        Congue possim facilisis sea justo dolore amet eos dolores est magna.`
    ],

    "grault": [
        // Date, Time, DateTime literals are available

        2020-12-31,              // -> time.Time 2020-12-31:00:00.000Z
        18:20:30.001,            // -> time.Time 1870-01-01T18:20:30.001Z
        2020-12-31T18:20:30.001Z // -> time.Time 2020-12-31:20:30.001Z
    ],

    // "key = value" syntax is allowed
    garply = true,

    // "key => value" syntax is allowed
    waldo => false,

    // Trailing commas are allowed.
    fred: 10,
}
```

### TOML

```toml
# Hash comment
// Line comment (non-standard)
/* Block comment (non-standard) */

"quoted-key" = "value0"
key          = "value1"
bare_key     = "value2"
bare-key     = "value3"  # Keys containing hyphens are allowed.
1234         = "value4"  # Keys starting with a number are allowed.
-3.14        = "value5"  # "-3": { "14": "value5" }
あいうえお    = "value6"  # // Non-ASCII identifiers are allowed. (non-standard)

[table-A]
item-a = 1
item-b = 2
item-c = "standard\n string"
item-d = '\literal string'
item-e = """\
         Multiline \
         standard string"""
item-f = '''
Multiline
literal string''' "\t" '''!!'''

[table-A.sub-table-P]
item-m = 11
item-n = 12
sub-table-X.sub-table-Y = 111

[[array-of-table-U]]
foo = 1
[array-of-table-U.sub-table-V]
bar = 11

[[array-of-table-U]]
foo = 2
[array-of-table-U.sub-table-V]
bar = 22
```

## 📚 Syntax

### Array

```js
[]
```

```js
[1]
```

```js
[1,]
```

```js
[1, 2]
```

```js
[1, 2,]
```

### Object

```js
{}
```

```js
{ "foo": 1 }
```

```js
{ "foo": 1, }
```

```js
{ 'foo': 1, }
```

```js
{ `foo`: 1, }
```

```js
{ foo: 1, }
```

```js
{ "foo": 1, "bar": 2, }
```

```js
{ "foo" => 1 }
```

```js
{ "foo" = 1 }
```

```js
{ foo.bar.baz = 1 }
// -> { "foo": { "bar": { "baz": 1 } } }
```

```js
{ "foo".bar."baz" = 1 }
// -> { "foo": { "bar": { "baz": 1 } } }
```

### Number

```js
123
```
```js
-123.45
```
```js
-1.2345e+6
```
```js
-123_456_789
```
```js
0x12345678
```
```js
0x1234_5678
```
```js
0o7654_3210
```
```js
0b0101_0101
```
```js
NaN
```
```js
Infinity
```
```js
+Infinity
```
```js
-Infinity
```

### String

```js
"foobar"
```
```js
'foobar'
```
```js
`foo
bar`
```
```js
"Hello\n\r\v\t\b\fWorld!"
```
```js
"\u0048\u0065\u006c\u006c\u006f\u002c\u0020\u0057\u006f\u0072\u006c\u0064\u0021"
```
```js
"\u{000048}\u{000065}\u{00006c}\u{00006c}\u{00006f}\u{00002c}\u{000020}\u{000057}\u{00006f}\u{000072}\u{00006c}\u{000064}\u{000021}"
```
```js
"\x48\x65\x6c\x6c\x6f\x2c\x20\x57\x6f\x72\x6c\x64\x21"
```

```js
"foo" "bar"
// -> "foobar"
```
```js
'foo' 'bar'
// -> 'foobar'
```

### Boolean

```js
true
```
```js
false
```

### Null, Undefined

```js
null
```
```js
undefined
```

### Comment

```sh
# Hash line comment
```
```js
// Line comment
```
```js
/*
   Block comment
*/
```

## ⚖️ License

MIT  
Copyright (c) 2023 Shellyl_N and Authors.
