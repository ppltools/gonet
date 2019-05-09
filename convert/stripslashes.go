package convert

import "unicode"

func StripcSlashes(data []byte) []byte {
	slen := len(data)

	ret := make([]byte, 0, slen)
	var tmp []byte

	for i := 0; i < slen; i++ {
		if data[i] == '\\' && i+1 < slen {
			i++
			switch data[i] {
			case 'n':
				ret = append(ret, '\n')
			case 'r':
				ret = append(ret, '\r')
			case 'a':
				ret = append(ret, '\a')
			case 't':
				ret = append(ret, '\t')
			case 'v':
				ret = append(ret, '\v')
			case 'b':
				ret = append(ret, '\b')
			case 'f':
				ret = append(ret, '\f')
			case '\\':
				ret = append(ret, '\\')
			case 'x':
				if i+1 < slen && isxdigit(data[i+1]) {
					tmp = append(tmp, data[i+1])
					i++
					if i+1 < slen && isxdigit(data[i+1]) {
						tmp = append(tmp, data[i+1])
						i++
					}
					ret = append(ret, strtol(tmp, 16))
				}
			default:
				j := 0
				for j+i < slen && data[j+i] >= '0' && data[j+i] <= '8' && j < 3 {
					tmp = append(tmp, data[i+j])
					j++
				}
				i = i + j
				if j != 0 {
					ret = append(ret, strtol(tmp, 8))
					i--
					tmp = tmp[:0]
				} else {
					ret = append(ret, data[i])
				}
			}
		} else {
			ret = append(ret, data[i])
		}
	}

	return ret
}

func isxdigit(b byte) bool {
	return unicode.IsDigit(rune(b)) || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func strtol(buf []byte, base int) byte {
	slen := len(buf)
	res := 0
	for i := 0; i < slen; i++ {
		var c int
		if unicode.IsDigit(rune(buf[i])) {
			c = int(buf[i] - '0')
		} else if buf[i] >= 'a' && buf[i] <= 'f' {
			c = int(buf[i]-'a') + 10
		} else if buf[i] >= 'A' && buf[i] <= 'F' {
			c = int(buf[i]-'A') + 10
		} else {
			return byte(res)
		}
		res *= base
		res += c
	}
	return byte(res)
}
