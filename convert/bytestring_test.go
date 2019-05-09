package convert

import (
	"bytes"
	"testing"
)

func TestBytesToString(t *testing.T) {
	data := map[string]struct {
		Test     []byte
		Expected string
	}{
		"t1": {[]byte("okkkk"), "okkkk"},
	}
	for name, cas := range data {
		res := BytesToString(cas.Test)
		if res != cas.Expected {
			t.Errorf("test %s failed: expected: %s, but is: %s", name, cas.Expected, res)
		}
	}
}

func TestStringToBytes(t *testing.T) {
	data := map[string]struct {
		Test     string
		Expected []byte
	}{
		"t1": {"okkkk", []byte("okkkk")},
	}
	for name, cas := range data {
		res := StringToBytes(cas.Test)
		if !bytes.Equal(res, cas.Expected) {
			t.Errorf("test %s failed: expected: %s, but is: %s", name, cas.Expected, res)
		}
	}
}
