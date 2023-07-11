package patb

// CharClass はキャラクタクラスを表します.
//
// r が適合するかを判定する関数です.
type CharClass func(r rune) bool

// C は set のいずれかの文字にマッチする CharClass を返します.
func C(set string) CharClass {
	v := []rune(set)
	l := len(v)
	if l <= 16 {
		return cfn[l](v)
	}
	return cmap(v)
}

// Not は set 以外の文字にマッチする CharClass を返します.
func Not(set string) CharClass {
	c := C(set)
	return func(r rune) bool {
		return !c(r)
	}
}

// Range は lo, hi の範囲にマッチする CharClass を返します.
//
// マッチする範囲に lo, hi を含みます.
func Range(lo, hi rune) CharClass {
	return func(r rune) bool {
		return lo <= r && r <= hi
	}
}

// All は全ての文字にマッチする CharClass を返します.
func All() CharClass {
	return func(r rune) bool {
		return true
	}
}

// Digit は 0-9 にマッチする CharClass を返します.
func Digit() CharClass {
	return func(r rune) bool {
		return '0' <= r && r <= '9'
	}
}

// Lower は英小文字 a-z にマッチする CharClass を返します.
func Lower() CharClass {
	return func(r rune) bool {
		return 'a' <= r && r <= 'z'
	}
}

// Upper は英大文字 A-Z にマッチする CharClass を返します.
func Upper() CharClass {
	return func(r rune) bool {
		return 'A' <= r && r <= 'A'
	}
}

// Alphabet は英字 A-Z, a-z にマッチする CharClass を返します.
func Alphabet() CharClass {
	return func(r rune) bool {
		if 'A' <= r && r <= 'A' {
			return true
		} else if 'a' <= r && r <= 'z' {
			return true
		}
		return false
	}
}

// Alnum は英数字 0-9, A-Z, a-z にマッチする CharClass を返します.
func Alnum() CharClass {
	return func(r rune) bool {
		if '0' <= r && r <= '9' {
			return true
		} else if 'A' <= r && r <= 'Z' {
			return true
		} else if 'a' <= r && r <= 'z' {
			return true
		}
		return false
	}
}

// Word は 0-9, A-Z, a-z, _ にマッチする CharClass を返します.
func Word() CharClass {
	return func(r rune) bool {
		if '0' <= r && r <= '9' {
			return true
		} else if 'A' <= r && r <= 'Z' {
			return true
		} else if 'a' <= r && r <= 'z' {
			return true
		} else if r == '_' {
			return true
		}
		return false
	}
}

// Blank は半角スペース, タブにマッチする CharClass を返します.
func Blank() CharClass {
	return func(r rune) bool {
		return r == ' ' || r == '\t'
	}
}

// Space は半角スペース文字にマッチする CharClass を返します.
func Space() CharClass {
	return func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f'
	}
}

//
// Optimized C functions
//

var cfn = []func(v []rune) CharClass{
	cfn00,
	cfn01, cfn02, cfn03, cfn04,
	cfn05, cfn06, cfn07, cfn08,
	cfn09, cfn10, cfn11, cfn12,
	cfn13, cfn14, cfn15, cfn16,
}

func cmap(v []rune) CharClass {
	m := make(map[rune]bool)
	for _, r := range v {
		m[r] = true
	}
	return func(r rune) bool {
		return m[r]
	}
}

func cfn00(v []rune) CharClass {
	return func(r rune) bool {
		return false
	}
}

func cfn01(v []rune) CharClass {
	c00 := v[0]
	return func(r rune) bool {
		return r == c00
	}
}

func cfn02(v []rune) CharClass {
	c00, c01 := v[0], v[1]
	return func(r rune) bool {
		return r == c00 || r == c01
	}
}

func cfn03(v []rune) CharClass {
	c00, c01, c02 := v[0], v[1], v[2]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02
	}
}

func cfn04(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03
	}
}

func cfn05(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04 := v[4]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04
	}
}

func cfn06(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05 := v[4], v[5]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05
	}
}

func cfn07(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06 := v[4], v[5], v[6]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06
	}
}

func cfn08(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07
	}
}

func cfn09(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08 := v[8]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08
	}
}

func cfn10(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09 := v[8], v[9]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09
	}
}

func cfn11(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10 := v[8], v[9], v[10]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10
	}
}

func cfn12(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10, c11 := v[8], v[9], v[10], v[11]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10 || r == c11
	}
}

func cfn13(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10, c11 := v[8], v[9], v[10], v[11]
	c12 := v[12]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10 || r == c11 ||
			r == c12
	}
}

func cfn14(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10, c11 := v[8], v[9], v[10], v[11]
	c12, c13 := v[12], v[13]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10 || r == c11 ||
			r == c12 || r == c13
	}
}

func cfn15(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10, c11 := v[8], v[9], v[10], v[11]
	c12, c13, c14 := v[12], v[13], v[14]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10 || r == c11 ||
			r == c12 || r == c13 || r == c14
	}
}

func cfn16(v []rune) CharClass {
	c00, c01, c02, c03 := v[0], v[1], v[2], v[3]
	c04, c05, c06, c07 := v[4], v[5], v[6], v[7]
	c08, c09, c10, c11 := v[8], v[9], v[10], v[11]
	c12, c13, c14, c15 := v[12], v[13], v[14], v[15]
	return func(r rune) bool {
		return r == c00 || r == c01 || r == c02 || r == c03 ||
			r == c04 || r == c05 || r == c06 || r == c07 ||
			r == c08 || r == c09 || r == c10 || r == c11 ||
			r == c12 || r == c13 || r == c14 || r == c15
	}
}
