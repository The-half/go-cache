package go_cache

import "runtime"

type Cache interface {
	Set(key string, value interface{})

	Get(key string) interface{}

	Del(key string)

	DelOldest()

	Len() int
}

type Value interface {
	Len() int
}

func CalcLen(value interface{}) int {
	var n int
	switch v := value.(type) {
	case Value:
		n = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 8 + len(v)
		}
	case int8, uint8, bool:
		n = 1
	case int16, uint16:
		n = 2
	case int32, uint32, float32:
		n = 4
	case int64, uint64, float64:
		n = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 4
		} else {
			n = 8
		}
	case complex64:
		n = 8
	case complex128:
		n = 16
	}

	return n
}

