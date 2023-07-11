package patb

import "testing"

func TestCharClass(t *testing.T) {
	tests := []struct {
		name string
		c    CharClass
		want map[rune]bool
	}{
		{`C("")`, C(""), map[rune]bool{
			'0': false, 'A': false, 'a': false, 'z': false,
			'あ': false, ' ': false, '\n': false,
		}},
		{`C("a")`, C("a"), map[rune]bool{
			'0': false, 'A': false, 'a': true, 'z': false,
			'あ': false, ' ': false, '\n': false,
		}},
		{`C("ab")`, C("ab"), map[rune]bool{
			'0': false, 'A': false, 'a': true, 'z': false,
			'あ': false, ' ': false, '\n': false,
		}},
		{`C("0Aaあ")`, C("0Aaあ"), map[rune]bool{
			'0': true, 'A': true, 'a': true, 'z': false,
			'あ': true, ' ': false, '\n': false,
		}},
		{`C("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")`, C("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"), map[rune]bool{
			'0': true, 'A': true, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Not("0Aaあ")`, Not("0Aaあ"), map[rune]bool{
			'0': false, 'A': false, 'a': false, 'z': true,
			'あ': false, ' ': true, '\n': true,
		}},
		{`Range('a', 'z')`, Range('a', 'z'), map[rune]bool{
			'0': false, 'A': false, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Dot()`, All(), map[rune]bool{
			'0': true, 'A': true, 'a': true, 'z': true,
			'あ': true, ' ': true, '\n': true,
		}},
		{`Digit()`, Digit(), map[rune]bool{
			'0': true, 'A': false, 'a': false, 'z': false,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Lower()`, Lower(), map[rune]bool{
			'0': false, 'A': false, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Upper()`, Upper(), map[rune]bool{
			'0': false, 'A': true, 'a': false, 'z': false,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Alphabet()`, Alphabet(), map[rune]bool{
			'0': false, 'A': true, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Alnum()`, Alnum(), map[rune]bool{
			'0': true, 'A': true, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Word()`, Word(), map[rune]bool{
			'0': true, 'A': true, 'a': true, 'z': true,
			'あ': false, ' ': false, '\n': false,
		}},
		{`Blank()`, Blank(), map[rune]bool{
			'0': false, 'A': false, 'a': false, 'z': false,
			'あ': false, ' ': true, '\n': false,
		}},
		{`Space()`, Space(), map[rune]bool{
			'0': false, 'A': false, 'a': false, 'z': false,
			'あ': false, ' ': true, '\n': true,
		}},
	}
	for _, te := range tests {
		for r, want := range te.want {
			got := te.c(r)
			if got != want {
				t.Errorf("%s ('%c') = %t, want %t", te.name, r, got, want)
			}
		}
	}
}
