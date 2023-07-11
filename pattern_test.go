package patb

import "testing"

func TestPattern(t *testing.T) {
	tests := []struct {
		name string
		pat  Pattern
		want map[string]int
	}{
		{`.`, Dot(), map[string]int{
			"aiu": 1,
			"あいう": 3,
		}},
		{`\w?`, Ch(0, 1, Word()), map[string]int{
			"aiu": 1,
			"あいう": 0,
		}},
		{`\w+`, Ch(1, 3, Word()), map[string]int{
			"aiu": 3,
			"あいう": -1,
		}},
		{`[\wあいう]+`, Ch(1, 3, Word(), C("あいう")), map[string]int{
			"aiu": 3,
			"あいう": 9,
		}},
		{`aiu`, S("aiu"), map[string]int{
			"aiu": 3,
			"あいう": -1,
		}},
		{`\w+@\w+`, Block(Ch(1, 3, Word()), S("@"), Ch(1, 3, Word())), map[string]int{
			"a@b":     3,
			"abc@xyz": 7,
			"abc":     -1,
			"@xyz":    -1,
		}},
		{`(\w+@)?`, Repeat(0, 1, Ch(1, 3, Word()), S("@")), map[string]int{
			"":     0,
			"abc":  0,
			"abc@": 4,
		}},
		{`(\w+\.)+`, Repeat(1, 2, Ch(1, 3, Word()), S(".")), map[string]int{
			"":            -1,
			"abc":         -1,
			"abc.":        4,
			"abc.xyz":     4,
			"abc.xyz.":    8,
			"abc.xyz.com": 8,
		}},
		{`(abc|xyz)`, Any(S("abc"), S("xyz")), map[string]int{
			"abc": 3,
			"xyz": 3,
			"aba": -1,
		}},
		{
			`^[^@]+@(\w+\.)+\w+$`,
			Block(
				Head(),
				Ch(1, 16, Not("@")),
				S("@"),
				Repeat(
					1, 3,
					Ch(1, 16, Word()),
					S("."),
				),
				Ch(1, 16, Word()),
				Tail(),
			),
			map[string]int{
				"a@b":           -1,
				"a@b.":          -1,
				"a@b.c":         5,
				"dum.my@go.dev": 13,
			},
		},
	}
	for _, te := range tests {
		for s, want := range te.want {
			got := te.pat(s, 0)
			if got != want {
				t.Errorf("%s ('%s') = %d, want %d", te.name, s, got, want)
			}
		}
	}
}
