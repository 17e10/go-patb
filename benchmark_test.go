package patb

import (
	"regexp"
	"testing"
)

func BenchmarkCharClass(b *testing.B) {
	texts := []string{
		`"display_name"<sip:0312341234@10.0.0.1:5060>;user=phone;hogehoge`,
		`<sip:0312341234@10.0.0.1>`,
		`"display_name"<sip:0312341234@10.0.0.1>`,
		`<sip:whois.this>;user=phone`,
		`"0333334444"<sip:[2001:30:fe::4:123]>;user=phone`,
	}

	tests := []struct {
		name  string
		class CharClass
	}{
		{"Word", Word()},
		{"C0", C("")},
		{"C1", C("a")},
		{"C2", C("ab")},
		{"C8", C("abcdefgh")},
		{"C16", C("abcdefghijklmnop")},
		{"Cx", C("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")},
	}
	for _, te := range tests {
		b.Run(te.name, func(b *testing.B) {
			class := te.class
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				for _, text := range texts {
					for _, c := range text {
						class(c)
					}
				}
			}
		})
	}
}

// BenchmarkCharClass/Word-4         	 3015433	       402.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/C0-4           	 3725443	       321.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/C1-4           	 2608927	       457.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/C2-4           	 2485777	       480.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/C8-4           	 2397547	       497.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/C16-4          	 1717198	       698.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCharClass/Cx-4           	  862430	      1337 ns/op	       0 B/op	       0 allocs/op

func BenchmarkSip(b *testing.B) {
	tests := []string{
		`"display_name"<sip:0312341234@10.0.0.1:5060>;user=phone;hogehoge`,
		`<sip:0312341234@10.0.0.1>`,
		`"display_name"<sip:0312341234@10.0.0.1>`,
		`<sip:whois.this>;user=phone`,
		`"0333334444"<sip:[2001:30:fe::4:123]>;user=phone`,
	}

	b.Run("regexp", func(b *testing.B) {
		re := regexp.MustCompile(`^["]{0,1}([^"]*)["]{0,1}[ ]*<(sip|tel|sips):(([^@]*)@){0,1}([^>^:]*|\[[a-fA-F0-9:]*\]):{0,1}([0-9]*){0,1}>(;.*){0,1}$`)
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for i, l := 0, len(tests); i < l; i++ {
				re.MatchString(tests[i])
			}
		}
	})
	b.Run("pat", func(b *testing.B) {
		pat := Block(
			Repeat(0, 1,
				S(`"`),
				Ch(1, 256, Not(`"`)),
				S(`"`),
				Ch(0, 256, C(" ")),
			),
			S("<"),
			Any(
				S("sip"),
				S("tel"),
				S("sips"),
			),
			S(":"),
			Repeat(0, 1,
				Ch(1, 256, Not("@")),
				S("@"),
			),
			Any(
				Block(
					S("["),
					Ch(1, 256, Alnum(), C(":")),
					S("]"),
				),
				Ch(0, 256, Not(">:")),
			),
			Repeat(0, 1,
				S(":"),
				Ch(1, 5, Digit()),
			),
			S(">"),
			Repeat(0, 1,
				S(";"),
				Ch(0, 256, All()),
			),
		)
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for i, l := 0, len(tests); i < l; i++ {
				pat(tests[i], 0)
			}
		}
	})
}

// BenchmarkSpeed/regexp-4         	  199252	      5916 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSpeed/pat-4         	  899504	      1263 ns/op	       0 B/op	       0 allocs/op
