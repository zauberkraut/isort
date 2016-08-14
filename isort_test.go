package isort

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

const (
	maxRandBits = 16
	maxRandLen  = 1<<maxRandBits - 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomize(a []int) {
	for i := range a {
		a[i] = rand.Int()
	}
}

/* This is a poor man's RNG of uniformly-distributed magnitudes. It should
   produce a much smaller mean than the 2**30 of a uniform 32-bit generator. */
func randMagOdd(nbits int) int {
	if nbits == 0 {
		return 0
	}
	if nbits > 32 {
		panic("requested random magnitude exceeds 32 bits")
	}

	mag := rand.Intn(nbits)
	return 1<<uint(mag) + randMagOdd(mag)
}

func randMag() int {
	n := randMagOdd(maxRandBits)
	return n ^ rand.Intn(2) // randomize parity
}

func benchSort(b *testing.B, sortf func([]int)) {
	var arr [maxRandLen]int

	b.StopTimer()
	b.ResetTimer()
	// TODO: this len randomization is somewhat stupid
	for i := 0; i < b.N; i++ {
		a := arr[:randMag()]
		randomize(a)
		b.StartTimer()
		sortf(a)
		b.StopTimer()
	}
}

func BenchmarkMsort(b *testing.B) {
	buf := make([]int, maxRandLen) // save on reallocation costs
	wrapper := func(a []int) {
		msort(a, buf)
	}
	benchSort(b, wrapper)
}

func BenchmarkQsortRandom(b *testing.B) {
	wrapper := func(a []int) {
		Qsort(a, Random)
	}
	benchSort(b, wrapper)
}

func BenchmarkQsortMiddle(b *testing.B) {
	wrapper := func(a []int) {
		Qsort(a, Middle)
	}
	benchSort(b, wrapper)
}

func BenchmarkQsortMedian3(b *testing.B) {
	wrapper := func(a []int) {
		Qsort(a, Median3)
	}
	benchSort(b, wrapper)
}

func BenchmarkQsortNinther(b *testing.B) {
	wrapper := func(a []int) {
		Qsort(a, Ninther)
	}
	benchSort(b, wrapper)
}

/* Benchmarks sort.Ints() */
func BenchmarkInts(b *testing.B) {
	benchSort(b, sort.Ints)
}

func TestMQsort(t *testing.T) {
	randUnsorted := make([]int, randMag())
	for i := range randUnsorted {
		randUnsorted[i] = rand.Int()
	}

	tests := [][]int{
		[]int{},
		[]int{0},
		[]int{0, 0},
		[]int{0, 0, 0},
		[]int{0, 1},
		[]int{1, 0},
		[]int{0, 0, 1},
		[]int{0, 1, 0},
		[]int{1, 0, 0},
		[]int{0, 1, 1},
		[]int{1, 0, 1},
		[]int{1, 1, 0},
		[]int{0, 1, 2},
		[]int{0, 2, 1},
		[]int{1, 0, 2},
		[]int{1, 2, 0},
		[]int{2, 0, 1},
		[]int{2, 1, 0},
		randUnsorted,
	}
	pivs := []Pivoter{Random, Middle, Median3, Ninther}

	for i, test := range tests {
		a := make([]int, len(test))
		copy(a, test)
		Msort(a)
		if !sort.IntsAreSorted(a) {
			t.Fatalf("Msort of test %d failed: %v => %v\n", i, test, a)
		}

		for j, piv := range pivs {
			copy(a, test)
			Qsort(a, piv)
			if !sort.IntsAreSorted(a) {
				t.Errorf("Qsort with pivoter %d of test %d failed: %v => %v\n", j, i, test, a)
			}
		}
	}
}
