# go-patb

[![GoDev][godev-image]][godev-url]

patb パッケージは正規表現 LIKE な文字列パターンマッチを提供します.

## Usage

```go
import "github.com/17e10/go-patb"

// `[^@]+@(\w+\.)+\w+`
pat := patb.Block(
    patb.Ch(1, 16, patb.Not("@")),
    patb.S("@"),
    patb.Repeat(
        1, 3,
        patb.Ch(1, 16, patb.Word()),
        patb.S("."),
    ),
    patb.Ch(1, 16, patb.Word()),
)
patb.Equal(pat, "a@b") == false
patb.Equal(pat, "a@b.c") == true
```

## License

This software is released under the MIT License, see LICENSE.

## Author

17e10

[godev-image]: https://pkg.go.dev/badge/github.com/17e10/go-patb
[godev-url]: https://pkg.go.dev/github.com/17e10/go-patb
