package patb

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/17e10/go-notifyb"
)

var SkipAll = notifyb.Notify("skip all")

type Writer interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

// Equal は s が pat と完全に一致するかを返します.
//
// pat に Head, Tail を含む必要はありません.
//
// 部分一致を判定する場合は Match を使用します.
// Equal は次の記述と機能は同じですが Equal の方が高速に動作します.
//
//	Match(Dot(), Block(Head(), pat, Tail()), s)
func Equal(pat Pattern, s string) bool {
	return pat(s, 0) == len(s)
}

// Match は s の中に pat と一致する部分があるかを返します.
func Match(c CharClass, pat Pattern, s string) bool {
	f, _ := FindIndex(c, pat, s, 0)
	return f >= 0
}

// FindIndex は s[i:] から pat に一致する範囲を返します.
// パターンに一致した文字列は s[f:l] です.
// 一致する部分がなければ -1, -1 を返します.
func FindIndex(c CharClass, pat Pattern, s string, i int) (f int, l int) {
	t := s[i:]
	for _, r := range t {
		if c(r) {
			if next := pat(s, i); next >= 0 {
				return i, next
			}
		}
		i += utf8.RuneLen(r)
	}
	return -1, -1
}

// FindAllFunc は s の中から pat と一致する部分を fn に渡します.
// fn が error を返すとそのエラーを返します.
func FindAllFunc(c CharClass, pat Pattern, s string, fn func(m string) error) error {
	var f, l int
	for {
		f, l = FindIndex(c, pat, s, f)
		if f < 0 {
			return nil
		}
		if err := fn(s[f:l]); err != nil {
			return err
		}
		f = l
	}
}

// ReplaceWrite は Writer を使って文字列を置換します.
//
// Writer を使用した文字列置換は頻繁なメモリアロケーションが発生せず柔軟に置換できるアプローチです.
// 置換後の文字列は Writer から取り出します.
//
// コールバック関数はパターンに一致した文字列を受け取り, 対応する置換文字列を Writer に書き込みます.
// コールバック関数が SkipAll を返すと, 以降の置換をスキップします.
// それ以外の error を返された場合, ReplaceWrite はその error を返します.
func ReplaceWrite(w Writer, c CharClass, pat Pattern, s string, fn func(w Writer, m string) error) error {
	i := 0
	for {
		f, l := FindIndex(c, pat, s, i)
		if f < 0 {
			w.WriteString(s[i:])
			break
		}
		if i < f {
			w.WriteString(s[i:f])
		}
		err := fn(w, s[f:l])
		if err == SkipAll {
			w.WriteString(s[l:])
			break
		} else if err != nil {
			return err
		}
		i = l
	}
	return nil
}

// ReplaceAll は src の pat に一致する部分をすべて repl に置き換えます.
func ReplaceAll(c CharClass, pat Pattern, src, repl string) string {
	w := bytes.NewBuffer(make([]byte, 0, len(src)))
	ReplaceWrite(w, c, pat, src, func(w Writer, _ string) error {
		w.WriteString(repl)
		return nil
	})
	return w.String()
}
