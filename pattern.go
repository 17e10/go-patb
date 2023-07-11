package patb

import (
	"unicode/utf8"
)

// Pattern はマッチする文字列パターンを表します.
//
// s[i:]の文字列が文字列パターンにマッチする関数です.
// マッチするとマッチした文字列の次の文字のインデックスを返します.
// マッチしなかった時は -1 を返します.
type Pattern func(s string, i int) int

// Dot は任意の 1 文字にマッチする Pattern です.
func Dot() Pattern {
	return func(s string, i int) int {
		if i > len(s) {
			return -1
		}
		_, w := utf8.DecodeRuneInString(s[i:])
		return i + w
	}
}

// Ch はキャラクタクラスにマッチする Pattern を返します.
//
// min, max の文字数の文字列にマッチします.
// キャラクタクラスは複数指定できます.
func Ch(min, max uint, classes ...CharClass) Pattern {
	var c CharClass
	switch len(classes) {
	case 0:
		c = func(r rune) bool {
			return false
		}
	case 1:
		c = classes[0]
	default:
		c = func(r rune) bool {
			for _, c := range classes {
				if c(r) {
					return true
				}
			}
			return false
		}
	}

	return func(s string, i int) int {
		if i > len(s) {
			return -1
		}
		s = s[i:]
		n := uint(0)
		for _, r := range s {
			if !c(r) {
				break
			}
			i, n = i+utf8.RuneLen(r), n+1
			if n >= max {
				break
			}
		}
		if n < min {
			return -1
		}
		return i
	}
}

// S は指定文字列にマッチする Pattern を返します.
func S(substr string) Pattern {
	w := len(substr)
	return func(s string, i int) int {
		if i+w > len(s) || s[i:i+w] != substr {
			return -1
		}
		return i + w
	}
}

// Head は先頭を表す Pattern を返します.
func Head() Pattern {
	return func(s string, i int) int {
		if i > 0 {
			return -1
		}
		return i
	}
}

// Tail は末尾を表す Pattern を返します.
func Tail() Pattern {
	return func(s string, i int) int {
		if i < len(s) {
			return -1
		}
		return i
	}
}

// Block は順次指定した Pattern とマッチする Pattern を返します.
func Block(pats ...Pattern) Pattern {
	return func(s string, i int) int {
		var next int
		for _, pat := range pats {
			if next = pat(s, i); next < 0 {
				return -1
			}
			i = next
		}
		return i
	}
}

// Repeat は min, max 回の繰り返しにマッチする Pattern を返します.
//
//	min, max の使用方法:
//		0, 1 ... 正規表現の ? と同等です.
//		0, n ... 正規表現の * に似た評価をします.
//		1, n ... 正規表現の + に似た評価をします.
func Repeat(min, max uint, pats ...Pattern) Pattern {
	pat := Block(pats...)
	return func(s string, i int) int {
		var next int
		var n uint
		for ; n < max; n++ {
			if next = pat(s, i); next < 0 {
				break
			}
			i = next
		}
		if n < min {
			return -1
		}
		return i
	}
}

// Any は指定した Pattern のいずれかとマッチする Pattern を返します.
func Any(pats ...Pattern) Pattern {
	return func(s string, i int) int {
		for _, pat := range pats {
			if next := pat(s, i); next >= 0 {
				return next
			}
		}
		return -1
	}
}
