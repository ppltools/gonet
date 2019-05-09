package convert

import (
	"unicode/utf8"

	"github.com/ppltools/utils/pool"
)

const lowerhex = "0123456789abcdef"

// Quote returns a double-quoted Go string literal representing s. The
// returned string uses Go escape sequences (\t, \n, \xFF, \u0100) for
// control characters and non-printable characters as defined by
// IsPrint.
func Quote(s string) (string, func()) {
	return quoteWith(s, '"', false, false)
}

func quoteWith(s string, quote byte, ASCIIonly, graphicOnly bool) (string, func()) {
	buf := pool.GetBuf(3*len(s)/2, true)
	new := appendQuotedWith(buf, s, quote, ASCIIonly, graphicOnly)
	return BytesToString(new), func() { pool.PutBuf(buf) }
}

func appendQuotedWith(buf []byte, s string, quote byte, ASCIIonly, graphicOnly bool) []byte {
	buf = append(buf, quote)
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[s[0]>>4])
			buf = append(buf, lowerhex[s[0]&0xF])
			continue
		}
		buf = appendEscapedRune(buf, r, quote, ASCIIonly, graphicOnly)
	}
	buf = append(buf, quote)
	return buf
}

func appendEscapedRune(buf []byte, r rune, quote byte, ASCIIonly, graphicOnly bool) []byte {
	var runeTmp [utf8.UTFMax]byte
	if r == rune(quote) || r == '\\' { // always backslashed
		buf = append(buf, '\\')
		buf = append(buf, byte(r))
		return buf
	}
	if ASCIIonly {
		if r < utf8.RuneSelf && IsPrint(r) {
			buf = append(buf, byte(r))
			return buf
		}
	} else if IsPrint(r) || graphicOnly && isInGraphicList(r) {
		n := utf8.EncodeRune(runeTmp[:], r)
		buf = append(buf, runeTmp[:n]...)
		return buf
	}
	switch r {
	case '\a':
		buf = append(buf, `\a`...)
	case '\b':
		buf = append(buf, `\b`...)
	case '\f':
		buf = append(buf, `\f`...)
	case '\n':
		buf = append(buf, `\n`...)
	case '\r':
		buf = append(buf, `\r`...)
	case '\t':
		buf = append(buf, `\t`...)
	case '\v':
		buf = append(buf, `\v`...)
	default:
		switch {
		case r < ' ':
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(r)>>4])
			buf = append(buf, lowerhex[byte(r)&0xF])
		case r > utf8.MaxRune:
			r = 0xFFFD
			fallthrough
		case r < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		}
	}
	return buf
}

// contains reports whether the string contains the byte c.
func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

// bsearch16 returns the smallest i such that a[i] >= x.
// If there is no such i, bsearch16 returns len(a).
func bsearch16(a []uint16, x uint16) int {
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// bsearch32 returns the smallest i such that a[i] >= x.
// If there is no such i, bsearch32 returns len(a).
func bsearch32(a []uint32, x uint32) int {
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// TODO: IsPrint is a local implementation of unicode.IsPrint, verified by the tests
// to give the same answer. It allows this package not to depend on unicode,
// and therefore not pull in all the Unicode tables. If the linker were better
// at tossing unused tables, we could get rid of this implementation.
// That would be nice.

// IsPrint reports whether the rune is defined as printable by Go, with
// the same definition as unicode.IsPrint: letters, numbers, punctuation,
// symbols and ASCII space.
func IsPrint(r rune) bool {
	// Fast check for Latin-1
	if r <= 0xFF {
		if 0x20 <= r && r <= 0x7E {
			// All the ASCII is printable from space through DEL-1.
			return true
		}
		if 0xA1 <= r && r <= 0xFF {
			// Similarly for ¡ through ÿ...
			return r != 0xAD // ...except for the bizarre soft hyphen.
		}
		return false
	}

	// Same algorithm, either on uint16 or uint32 value.
	// First, find first i such that isPrint[i] >= x.
	// This is the index of either the start or end of a pair that might span x.
	// The start is even (isPrint[i&^1]) and the end is odd (isPrint[i|1]).
	// If we find x in a range, make sure x is not in isNotPrint list.

	if 0 <= r && r < 1<<16 {
		rr, isPrint, isNotPrint := uint16(r), isPrint16, isNotPrint16
		i := bsearch16(isPrint, rr)
		if i >= len(isPrint) || rr < isPrint[i&^1] || isPrint[i|1] < rr {
			return false
		}
		j := bsearch16(isNotPrint, rr)
		return j >= len(isNotPrint) || isNotPrint[j] != rr
	}

	rr, isPrint, isNotPrint := uint32(r), isPrint32, isNotPrint32
	i := bsearch32(isPrint, rr)
	if i >= len(isPrint) || rr < isPrint[i&^1] || isPrint[i|1] < rr {
		return false
	}
	if r >= 0x20000 {
		return true
	}
	r -= 0x10000
	j := bsearch16(isNotPrint, uint16(r))
	return j >= len(isNotPrint) || isNotPrint[j] != uint16(r)
}

// IsGraphic reports whether the rune is defined as a Graphic by Unicode. Such
// characters include letters, marks, numbers, punctuation, symbols, and
// spaces, from categories L, M, N, P, S, and Zs.
func IsGraphic(r rune) bool {
	if IsPrint(r) {
		return true
	}
	return isInGraphicList(r)
}

// isInGraphicList reports whether the rune is in the isGraphic list. This separation
// from IsGraphic allows quoteWith to avoid two calls to IsPrint.
// Should be called only if IsPrint fails.
func isInGraphicList(r rune) bool {
	// We know r must fit in 16 bits - see makeisprint.go.
	if r > 0xFFFF {
		return false
	}
	rr := uint16(r)
	i := bsearch16(isGraphic, rr)
	return i < len(isGraphic) && rr == isGraphic[i]
}
