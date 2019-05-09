package convert

import (
	"bytes"
	"testing"
)

func TestStripcSlashes(t *testing.T) {
	data := map[string]struct {
		Test     []byte
		Expected []byte
	}{
		"t1": {[]byte("How\\'s everybody"), []byte("How's everybody")},
		"t2": {[]byte("Are you \\\"JOHN\\\"?"), []byte("Are you \"JOHN\"?")},
		"t3": {[]byte("c:\\\\php\\\\stripslashes"), []byte("c:\\php\\stripslashes")},
		"t4": {[]byte("c:\\php\\stripslashes"), []byte("c:phpstripslashes")},
		"t5": {[]byte("\\xff"), []byte{255}},
	}
	for name, cas := range data {
		res := StripcSlashes(cas.Test)
		if !bytes.Equal(res, cas.Expected) {
			t.Errorf("test %s failed: expected: %s, but is: %s", name, cas.Expected, res)
		}
	}
}
