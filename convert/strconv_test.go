package convert

import (
	"testing"
)

func TestQuote(t *testing.T) {
	data := map[string]struct {
		Test     string
		Expected string
	}{
		"t1": {"ab", "\"ab\""},
	}
	for name, cas := range data {
		res, cancel := Quote(cas.Test)
		if res != cas.Expected {
			t.Fatalf("case %s failed: expected: %s, but got: %s", name, cas.Expected, res)
		}
		cancel()
	}
}

func BenchmarkQuote(b *testing.B) {
	b.ReportAllocs()
	expected := "\"ab\""
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf, cancel := Quote("ab")
			// do something
			if buf != expected {
				b.Fatalf("expected: %s, but got: %s", expected, buf)
			}
			cancel()
		}
	})
}
