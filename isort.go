/* isort.go

   These are toy implementations of mergesort and quicksort that operate only
   over int slices. They can substantially outperform sort.Ints(), but certainly
   aren't meant as replacements for the generalized sort package. Several
   pivoting strategies are provided; clients may also implement their own. */

package isort

import (
	"math"
	"math/rand"
)

func msort(a []int, b []int, depth int) {
	n := len(a)
	mid := n / 2

	if depth > 0 {
		if mid > 1 {
			msort(b[:mid], a[:mid], depth-1)
		}
		if n-mid > 1 {
			msort(b[mid:], a[mid:], depth-1)
		}

		for i, l, r := 0, 0, mid; i < n; i++ {
			if r == n || (l < mid && b[l] <= b[r]) {
				a[i] = b[l]
				l++
			} else {
				a[i] = b[r]
				r++
			}
		}
	} else if a[0] > a[1] { // in-place sorting at pair leaf
		a[0], a[1] = a[1], a[0]
	}
}

/* Top-down mergesort. */
func Msort(a []int) {
	if len(a) > 1 {
		b := make([]int, len(a))
		depth := int(math.Floor(math.Log2(float64(len(a)))))
		sortToB := depth&1 == 1
		if sortToB {
			msort(b, a, depth)
			copy(a, b)
		} else {
			msort(a, b, depth)
		}
	}
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
