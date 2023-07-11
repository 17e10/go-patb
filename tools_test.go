package patb

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestEqual(t *testing.T) {
	name := "[[:alpha:]]{3,5}"
	al35 := Ch(3, 5, Alphabet())

	tests := []struct {
		name string
		pat  Pattern
		s    string
		want bool
	}{
		{name, al35, "abc", true},
		{name, al35, "abcde", true},
		{name, al35, "abcdef", false},
	}
	for _, te := range tests {
		got := Equal(te.pat, te.s)
		if got != te.want {
			t.Errorf("Equal(`%s`, %q) = %t, want %t", te.name, te.s, got, te.want)
		}
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name string
		c    CharClass
		pat  Pattern
		s    string
		want bool
	}{
		{`\w{3,5}`, Word(), Ch(3, 5, Word()), "012abcあいう", true},
		{`\s{1,16}`, Space(), Ch(1, 16, Space()), "abc あいう", true},
		{`\s{1,16}`, Space(), Ch(1, 16, Space()), "abcあいう", false},
	}
	for _, te := range tests {
		got := Match(te.c, te.pat, te.s)
		if got != te.want {
			t.Errorf("Match(`%s`, %q) = %t, want %t", te.name, te.s, got, te.want)
		}
	}
}

func TestFindIndex(t *testing.T) {
	type rval struct {
		f int
		l int
	}
	tests := []struct {
		name string
		c    CharClass
		pat  Pattern
		s    string
		want rval
	}{
		{`\w{3,5}`, Word(), Ch(3, 5, Word()), "012abcあいう", rval{0, 5}},
		{`\s{1,16}`, Space(), Ch(1, 16, Space()), "abc あいう", rval{3, 4}},
		{`\s{1,16}`, Space(), Ch(1, 16, Space()), "abcあいう", rval{-1, -1}},
	}
	for _, te := range tests {
		var got rval
		got.f, got.l = FindIndex(te.c, te.pat, te.s, 0)
		if !reflect.DeepEqual(got, te.want) {
			t.Errorf("FindIndex(`%s`, %q) = %v, want %v", te.name, te.s, got, te.want)
		}
	}
}

func TestFindAllFunc(t *testing.T) {
	tests := []struct {
		name string
		c    CharClass
		pat  Pattern
		s    string
		want []string
	}{
		{`[^ ]{1,16}`, Not(" "), Ch(1, 16, Not(" ")), "abc あいう", []string{"abc", "あいう"}},
	}
	var got []string
	for _, te := range tests {
		got = got[:0]
		FindAllFunc(te.c, te.pat, te.s, func(m string) error {
			got = append(got, m)
			return nil
		})
		if !reflect.DeepEqual(got, te.want) {
			t.Errorf("FindAllFunc(`%s`, %q) = %v, want %v", te.name, te.s, got, te.want)
		}
	}
}

func TestReplaceWrite(t *testing.T) {
	tests := []struct {
		name string
		c    CharClass
		pat  Pattern
		s    string
		fn   func(w Writer, m string) error
		want string
	}{
		{
			"[^aiueo]",
			Not("aiueo"),
			Dot(),
			"seafood fool",
			func(w Writer, m string) error {
				w.WriteString(strings.ToUpper(m))
				return nil
			},
			"SeaFooD FooL",
		},
		{
			"[^aiueo]",
			Not("aiueo"),
			Dot(),
			"seafood fool",
			func(w Writer, m string) error {
				w.WriteString(strings.ToUpper(m))
				return SkipAll
			},
			"Seafood fool",
		},
	}
	w := bytes.NewBuffer(make([]byte, 0, 32))
	for _, te := range tests {
		w.Reset()
		err := ReplaceWrite(w, te.c, te.pat, te.s, te.fn)
		if err != nil {
			t.Errorf("%s errored %v", te.name, err)
			continue
		}
		got := w.String()
		if got != te.want {
			t.Errorf("%q replaced %q, want %q", te.s, got, te.want)
		}
	}
}

func TestSqlbT(t *testing.T) {
	c := C("@#$=!")
	pat := Any(
		Block(S("@"), Ch(0, 2, Digit())),
		Block(S("#"), Ch(0, 2, Digit())),
		Block(S("$"), Ch(0, 2, Digit())),
		Block(Ch(0, 1, C("!")), S("=="), Ch(1, 16, Space()), S("@"), Ch(0, 2, Digit())),
	)

	tests := []struct {
		s    string
		want []string
	}{
		{"a $ b", []string{"$"}},
		{"`age` = @", []string{"@"}},
		{"# = @", []string{"#", "@"}},
		{"# == @", []string{"#", "== @"}},
		{"# !== @", []string{"#", "!== @"}},
		{"#0 <= @1 AND #0 < @2", []string{"#0", "@1", "#0", "@2"}},
		{"#1 = @", []string{"#1", "@"}},
		{"INSERT INTO person (#) VALUES (@)", []string{"#", "@"}},
	}

	for _, te := range tests {
		var got []string
		FindAllFunc(c, pat, te.s, func(m string) error {
			got = append(got, m)
			return nil
		})
		if !reflect.DeepEqual(got, te.want) {
			t.Errorf("sql.T(%q) = %v, want %v", te.s, got, te.want)
		}
	}
}

func TestPostalBanchiKanji(t *testing.T) {
	banchiClass := func(r rune) bool {
		if '０' <= r && r <= '９' {
			return true
		} else if r == '～' || r == '－' || r == 'の' || r == '・' {
			return true
		}
		return false
	}

	c := C("無番大０１２３４５６７８９")
	pat := Any(
		Block(Head(), S("無番地を除く"), Tail()),
		Block(Head(), S("番地のみ"), Tail()),
		Block(Head(), S("大字"), Tail()),
		Block(Head(), S("番地"), Tail()),
		Block(Head(), Ch(1, 20, banchiClass), Tail()),
		S("無番地"),
		Block(
			Ch(1, 20, banchiClass),
			Any(
				S("線"),
				S("丁目"),
			),
		),
		Repeat(1, 16,
			Ch(1, 20, banchiClass),
			Repeat(0, 1,
				S("番地"),
				Any(
					S("～"),
					S("以上"),
					S("以降"),
					S(""),
				),
			),
			Ch(0, 1, Space()),
		),
	)

	tests := []struct {
		s    string
		want string
	}{
		{"番地", "番地"},
		{"番地のみ", "番地のみ"},
		{"無番地を除く", "無番地を除く"},
		{"大字", "大字"},
		{"１～１３１番地", "１～１３１番地"},
		{"４００", "４００"},
		{"４００－２番地", "４００－２番地"},
		{"西５～８線７９～１１０番地", "７９～１１０番地"},
		{"８０６番地", "８０６番地"},
		{"１丁目", ""},
		{"３丁目５", "５"},
		{"１３－４", "１３－４"},
		{"油駒", ""},
		{"１３２～１５６", "１３２～１５６"},
		{"４丁目５５～１１４番地", "５５～１１４番地"},
		{"４０の１番地", "４０の１番地"},
		{"新田１７－２", "１７－２"},
		{"３７番地", "３７番地"},
		{"東火行１番地", "１番地"},
		{"５３の１～６０の９番地", "５３の１～６０の９番地"},
		{"１の２", "１の２"},
		{"３の２～６", "３の２～６"},
		{"４の２・４・６", "４の２・４・６"},
		{"１１の１番地", "１１の１番地"},
		{"８９７番地", "８９７番地"},
		{"中島５０５～５１８番地", "５０５～５１８番地"},
		{"１７００番地～", "１７００番地～"},
		{"１～４２６番地（川東）", "１～４２６番地"},
		{"４２７番地以降（川西）", "４２７番地以降"},
		{"１１３～７９１番地", "１１３～７９１番地"},
		{"稲崎平３０２番地・３１５番地", "３０２番地・３１５番地"},
		{"南原無番地", ""},
		{"４３０番地以上", "４３０番地以上"},
		{"１～５００ 古町", "１～５００ "},
		{"１７３～２５７番地 鉢伏峠", "１７３～２５７番地 "},
	}
	for _, te := range tests {
		var got string
		FindAllFunc(c, pat, te.s, func(m string) error {
			if m == "無番地" {
				return nil
			}
			if strings.HasSuffix(m, "線") {
				return nil
			}
			if strings.HasSuffix(m, "丁目") {
				return nil
			}
			got = m
			return SkipAll
		})
		if got != te.want {
			t.Errorf("%q Matched %q, want %q", te.s, got, te.want)
		}
	}
}
