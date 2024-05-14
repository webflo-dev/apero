package utils

func Default[T any]() (ret T) {
	return
}

type Threshold[K int, V any] map[K]V

func (t Threshold[K, V]) GetThreshold(value K) (V, bool) {
	for k, v := range t {
		if value >= k {
			return v, true
		}
	}
	return Default[V](), false
}
