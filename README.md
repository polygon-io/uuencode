# Package UUEncode

A short and sweet Go library that supports decoding uuencoded things.

For more information on what uuencoding is/how it works, check out [this wikipedia article](https://en.wikipedia.org/wiki/Uuencoding).

**Important Note:** This package currently only supports _decoding_ uuencoded contents (because...well...that's all we need here at Polygon.io for now :shrug:).
Contributions are welcome, if you'd like to implement an `Encoder` struct and create a PR we'd be overjoyed :D

uuencoding is an old, rarely unused format at this point and the standard isn't very strict.
There are lots of little variations in different implementations. 

This particular implementation is geared towards decoding binary files within SEC filings. 
It implements the behavior described in the wikipedia article linked, so it should be relatively portable.
This implementation also adds some extra features to clean up input that doesn't quite conform to the expectations of that format.

There are tests ensuring this package works decoding standard input, input encoded via the `uuencode` utility on macOS, and input encoded in the style that the SEC follows. 

## Examples

For examples, check out the test files ([decode](./decode_test.go))