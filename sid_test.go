package sid

import (
	"fmt"
	"testing"
)

func Example() {
	e := New(4, 24, '0')

	fmt.Println(e.shuffle(12))
	fmt.Println(e.unshuffle(3145728))

	fmt.Println(e.Encode(12))
	fmt.Println(e.Decode("yNvD"))

	// Output:
	// 3145728
	// 12
	// yNvD
	// 12
}

func TestEncoder(t *testing.T) {
	e := New(4, 24, '0')
	for i := int64(0); i < 10000000; i++ {
		j := e.Decode(e.Encode(i))
		if i != j {
			t.Fatalf("want %d, got %d", i, j)
		}
	}
}

func BenchmarkEncoder(b *testing.B) {
	e := New(4, 24, '0')

	b.Run("encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.Encode(12)
		}
	})

	b.Run("decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.Decode("yNvD")
		}
	})
}
