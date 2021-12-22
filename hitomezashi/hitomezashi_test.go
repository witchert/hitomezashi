package hitomezashi

import (
	"testing"
)

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(
			[]bool{true, true, false, false, true, true, true, true, false, false, true, true, false, true, true, false, false, true, false, false, true, true, true, false},
			[]bool{true, true, false, false, true, true, true, true, false, false, true, true, false, true, true, false, false, true, false, false, true, true, true, false},
			24,
		)
	}
}
