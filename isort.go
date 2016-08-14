/* isort.go

   These are toy implementations of mergesort and quicksort that operate only
   over int slices. They can substantially outperform sort.Ints(), but certainly
   aren't meant as replacements for the generalized sort package. Several
   pivoting strategies are provided; clients may also implement their own. */

package isort

import (
	"math/rand"
)

func msort(a []int, b []int) {
	mid := len(a) / 2
	if mid >= 2 {
		msort(a[:mid], b)
	}
	if len(a)-mid >= 2 {
		msort(a[mid:], b)
	}

	copy(b, a)

	for i, j, k := 0, mid, 0; k < len(a); k++ {
		if j == len(a) || (i < mid && b[i] < b[j]) {
			a[k] = b[i]
			i++
		} else {
			a[k] = b[j]
			j++
		}
	}
}

/* Mergesort */
func Msort(a []int) {
	b := make([]int, len(a))
	msort(a, b)
}

/* Chooses a pivot and returns its (index, value). */
type Pivoter func([]int) (int, int)

func Random(a []int) (int, int) {
	i := rand.Intn(len(a))
	return i, a[i]
}

func Middle(a []int) (int, int) {
	i := len(a) / 2
	return i, a[i]
}

/* Selects the median of the first, middle and last elements of the slice. */
func Median3(a []int) (int, int) {
	i, j := len(a)/2, len(a)-1
	lo, mid, hi := a[0], a[i], a[j]
	switch {
	case mid <= lo && lo <= hi:
		return 0, lo
	case lo <= mid && mid <= hi:
		return i, mid
	default:
		return j, hi
	}
}

/* Selects the median among the medians-of-three of the thirds of the slice. */
func Ninther(a []int) (int, int) {
	if len(a) < 3 {
		return Median3(a)
	}

	midStart, midEnd := len(a)/3, 2*len(a)/3
	var i, m [3]int
	i[0], m[0] = Median3(a[:midStart])
	i[1], m[1] = Median3(a[midStart:midEnd])
	i[2], m[2] = Median3(a[midEnd:])
	i[1] += midStart
	i[2] += midEnd
	j, p := Median3(m[:])
	return i[j], p
}

/* Quicksort */
func Qsort(a []int, piv Pivoter) {
	if len(a) < 2 {
		return
	}

	k, p := piv(a)
	r := len(a) - 1
	a[k], a[r] = a[r], a[k]

	i := 0
	for j, x := range a[:r] {
		if x <= p {
			a[i], a[j] = x, a[i]
			i++
		}
	}

	a[i], a[r] = a[r], a[i]

	Qsort(a[:i], piv)
	Qsort(a[i+1:], piv)
}

/* A sensible default. */
func Sort(a []int) {
	Qsort(a, Median3)
}
