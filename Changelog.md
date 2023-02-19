# Changelog

# v0.0.15
* Using constants for ClassName.
* Edit README.

# v0.0.14
* [FIX] Fix unmarshaller; Conversion from map to struct causes an error if there is no key matching the field name in the map.

# v0.0.13
* Added unmarshaller.
* [Breaking change] The required minimum version of Go has changed from 1.17 to 1.18.
  * Change `go.mod` and `test.yml`.
* Edit README.
* Edit wasm example html.

# v0.0.12
* [Breaking change] Change interop parameter type.
* Edit README.

# v0.0.11
* Added complex number parser.
* Remove unnecessary `erase()` call.
* Added dummy `version.go` file: for calling `go get -u`.
* Edit README.

# v0.0.10
* Improve word boundary.
* Edit wasm example html.
* Edit README.

# v0.0.9
* Improve JSON string parser: support TOML style multiline strings.
* Edit wasm example.
* Edit README.

# v0.0.8
* Added hexadecimal float literals.
* Added tests.
* Edit wasm example html.

# v0.0.7
* Fix string parsers: escape sequence.
* Added strict IEEE-754 +0/-0 parsing test.
* Edit README.

# v0.0.6
* Added TOML parser.
* Edit README.
* Rename files.

# v0.0.5
* Added dotted key syntax.

# v0.0.4
* Added '=' to the key/value separator in maps.

# v0.0.3
* Add wasm example.

# v0.0.2
* (Breaking change) Change `Parse` parameters.
  * Add `interop` parameter.
* Improve error reporting.
* ~~Add wasm example.~~
* Edit README.

# v0.0.1
* First release.
