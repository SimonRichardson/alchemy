package hashring

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"github.com/spaolacci/murmur3"
)

func TestHashRingAddRemove(t *testing.T) {
	t.Parallel()

	t.Run("add", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 2)
			return ring.Add(a.String())
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("add duplicate", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 2)
			ring.Add(a.String())
			return !ring.Add(a.String())
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 2)
			return !ring.Remove(a.String())
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("add then remove", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 2)
			ring.Add(a.String())
			return ring.Remove(a.String())
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestHashRingLookup(t *testing.T) {
	t.Parallel()

	t.Run("lookup", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 10)
			if expected, actual := true, ring.Add(a.String()); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			want := []string{
				a.String(),
			}
			got := ring.LookupN(a.String(), 2)

			return reflect.DeepEqual(want, got)
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("lookup with larger corpus", func(t *testing.T) {
		fn := func(a []ASCII) bool {
			if len(a) < 2 {
				return true
			}

			ring := New(murmur3.Sum32, 10)
			for _, v := range a {
				ring.Add(v.String())
			}

			var (
				key = a[0].String()
				got = ring.LookupN(key, 2)
			)

			return len(got) == 2
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("lookup with empty value", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 10)

			want := []string{}
			got := ring.LookupN(a.String(), 2)

			return reflect.DeepEqual(want, got)
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestHashRingChecksum(t *testing.T) {
	t.Parallel()

	t.Run("add", func(t *testing.T) {
		fn := func(a ASCII) bool {
			ring := New(murmur3.Sum32, 64)
			ring.Add(a.String())

			v0, err0 := ring.Checksum()
			if err0 != nil {
				t.Fatal(err0)
			}

			v1, err1 := ring.Checksum()
			if err1 != nil {
				t.Fatal(err1)
			}

			return v0 == v1
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		fn := func(a, b ASCII) bool {
			ring := New(murmur3.Sum32, 64)
			ring.Add(a.String())
			ring.Add(a.String())
			ring.Add(b.String())
			ring.Remove(a.String())

			v0, err0 := ring.Checksum()
			if err0 != nil {
				t.Fatal(err0)
			}

			v1, err1 := ring.Checksum()
			if err1 != nil {
				t.Fatal(err1)
			}

			return v0 == v1
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("add and remove", func(t *testing.T) {
		fn := func(a, b ASCII) bool {
			ring := New(murmur3.Sum32, 64)
			ring.Add(a.String())
			ring.Add(a.String())
			ring.Add(b.String())

			v0, err0 := ring.Checksum()
			if err0 != nil {
				t.Fatal(err0)
			}

			ring.Remove(a.String())

			v1, err1 := ring.Checksum()
			if err1 != nil {
				t.Fatal(err1)
			}

			return v0 != v1
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

const asciiChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateString creates a random string
func GenerateString(r *rand.Rand, size int) string {
	res := make([]byte, size)
	for k := range res {
		res[k] = byte(asciiChars[r.Intn(len(asciiChars)-1)])
	}
	return string(res)
}

// ASCII creates a value that is simple ascii characters from a-Z0-9.
type ASCII string

// Generate allows ASCII to be used within quickcheck scenarios.
func (ASCII) Generate(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(ASCII(GenerateString(r, size)))
}

func (a ASCII) String() string {
	return string(a)
}

// ASCIISlice creates a series of values that are simple ascii characters from
// a-Z0-9
type ASCIISlice []string

// Generate allows ASCIISlice to be used within quickcheck scenarios.
func (ASCIISlice) Generate(r *rand.Rand, size int) reflect.Value {
	res := make([]string, size)
	for k := range res {
		res[k] = GenerateString(r, size)
	}
	return reflect.ValueOf(res)
}

// Slice returns the underlying slice of the type
func (a ASCIISlice) Slice() []string {
	return []string(a)
}

func (a ASCIISlice) String() string {
	return strings.Join(a.Slice(), ",")
}
