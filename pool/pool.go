package pool

import (
	"sync"
)

const (
	_NumSizeClasses = 14
)

var class_to_size = [_NumSizeClasses]int{32, 64, 128, 256, 512, 1024, 2048, 4096, 10240, 20480, 40960, 81920, 163840, 327680}

func getPoolSizeClass(size int) int {
	for i := 0; i < _NumSizeClasses; i++ {
		if class_to_size[i] >= size {
			return i
		}
	}
	return _NumSizeClasses - 1
}

var pools = []*sync.Pool{}

func init() {
	pools = make([]*sync.Pool, _NumSizeClasses)
	for i := 0; i < _NumSizeClasses; i++ {
		tmpSize := class_to_size[i]
		pools[i] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, tmpSize)
			},
		}
	}
}

// GetBuf: get one buffer from the proper pool
func GetBuf(size int, zero bool) []byte {
	class := getPoolSizeClass(size)
	p := pools[class].Get()
	buf := p.([]byte)
	if zero {
		return buf[:0]
	}
	return buf
}

// PutBuf: put the buffer back to the proper pool
func PutBuf(buf []byte) {
	class := getPoolSizeClass(len(buf))
	pools[class].Put(buf)
}
