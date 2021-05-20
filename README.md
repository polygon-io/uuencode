# Package UUEncode

A short and sweet Go library that supports decoding uuencoded things.

For more information on what uuencoding is/how it works, check out [this wikipedia article](https://en.wikipedia.org/wiki/Uuencoding).

**Important Note:** This package currently only supports _decoding_ uuencoded contents (because...well...that's all we need here at Polygon.io for now :shrug:).
Contributions are welcome, if you'd like to implement an `Encoder` struct and create a PR we'd be overjoyed :D

## Examples

For examples, check out the test files ([decode](./decode_test.go))