package sid

import (
	"bytes"
	"sync"
	"sync/atomic"
)

const (
	Alphabet          = "JedR8LNFY2j6MrhkBSADUyfP5amuH9xQCX4VqbgpsGtnW7vc3TwKE"
	AlphabetCanonical = "mn6j2c4rv8bpygw95z7hsdaetxuk3fq"
)

var bufferPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}

type Encoder struct {
	Alphabet string
	mask     int64
	pad      byte
	n        int
	b        int
	indices  [75]int
}

// New creates an sid encoder. n is the minimum length of
// the encoded ID. b specifies how many bits will be shuffled.
// any bits higher than b will remain as is. b of 0 will leave
// all bits unaffected and the algorithm will simply be converting
// your integer to a different base. pad is used when the length
// of encoded ID is less than n.
func New(n, b int, pad byte) *Encoder {
	e := Encoder{
		Alphabet: Alphabet,
		mask:     (1 << b) - 1,
		pad:      pad,
		n:        n,
		b:        b,
	}
	for i := 0; i < len(e.Alphabet); i++ {
		e.indices[e.Alphabet[i]-0x30] = i
	}
	return &e
}

// Encode encodes an integer to a sid.
func (e Encoder) Encode(x int64) string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	e.enbase(e.shuffle(x), buf)

	b := buf.Bytes()
	n := len(b)
	for i := 0; i < n/2; i++ {
		j := n - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// Decode decodes a sid back to an integer.
func (e Encoder) Decode(s string) int64 {
	return e.unshuffle(e.debase(s))
}

// Shuffle converts an integer to a shuffled integer.
func (e Encoder) Shuffle(x int64) int64 {
	return e.shuffle(x)
}

// Unshuffle converts an integer back to its original value.
func (e Encoder) Unshuffle(x int64) int64 {
	return e.unshuffle(x)
}

func (e Encoder) shuffle(x int64) int64 {
	var r int64
	y := x & e.mask
	for i := 0; i < e.b; i++ {
		if y&(1<<i) != 0 {
			r |= (1 << (e.b - i - 1))
		}
	}
	return r | (x &^ e.mask)
}

func (e Encoder) unshuffle(x int64) int64 {
	var r int64
	y := x & e.mask
	for i := 0; i < e.b; i++ {
		if y&(1<<(e.b-i-1)) != 0 {
			r |= (1 << i)
		}
	}
	return r | (x &^ e.mask)
}

func (e Encoder) enbase(x int64, buf *bytes.Buffer) {
	n := int64(len(e.Alphabet))

	for x >= n {
		buf.WriteByte(e.Alphabet[x%n])
		x /= n
	}

	buf.WriteByte(e.Alphabet[x%n])

	for i := 0; i < e.n-buf.Len(); i++ {
		buf.WriteByte(e.pad)
	}
}

func (e Encoder) debase(s string) int64 {
	var r int64
	n := len(e.Alphabet)
	k := len(s)
	for i := 0; i < k; i++ {
		r += int64(e.indices[s[k-i-1]-0x30] * pow(n, i))
	}
	return r
}

func pow(x, y int) int {
	if y == 0 {
		return 1
	}

	r := 1
	for y > 0 {
		if y%2 == 1 {
			r *= x
		}
		x *= x
		y /= 2
	}
	return r
}

var defaultEncoder atomic.Value

// Encode encodes an integer to a sid.
func Encode(x int64, n int) string {
	return defaultEncoder.Load().(*Encoder).Encode(x)
}

// Decode decodes a sid back to an integer.
func Decode(s string) int64 {
	return defaultEncoder.Load().(*Encoder).Decode(s)
}

// Shuffle converts an integer to a shuffled integer.
func Shuffle(x int64) int64 {
	return defaultEncoder.Load().(*Encoder).Shuffle(x)
}

// Unshuffle converts an integer back to its original value.
func Unshuffle(x int64) int64 {
	return defaultEncoder.Load().(*Encoder).Unshuffle(x)
}

func init() {
	defaultEncoder.Store(New(4, 24, '0'))
}
