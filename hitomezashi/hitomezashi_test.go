package hitomezashi

import (
	"testing"
)

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New([]string{"0", "1"}, []string{"1", "0"}, 24)
	}
}
