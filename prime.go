// Package prime provides functionality to produce prime numbers using all
// available cpu cores. https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes
// can be an starting point to find more information about how to calculate
// prime numbers.
//
// The method used in Primes function is more memory consuming than a simpler
// Trial division method (https://en.wikipedia.org/wiki/Trial_division)
package prime

import (
	"math"
	"runtime"
)

func fill(nums []bool, i uint64, max uint64) {
	a := 2 * i
	for a <= max {
		nums[a] = true
		a = a + i
	}
}

func goFill(nums []bool, i uint64, max uint64, next <-chan bool) {
	fill(nums, i, max)
	<-next
}

// Primes returns a slice of all prime numbers equal or lower than max.
func Primes(max uint64) []uint64 {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	next := make(chan bool, cores)

	var nums = make([]bool, max+1)
	m := uint64(math.Sqrt(float64(max)))

	for i := uint64(2); i < m; i = i + 2 {
		if nums[i] == false {
			go goFill(nums, i, max, next)
			next <- true
		}
		if i == 2 {
			i = 1
		}
	}

	for i := 0; i < cores; i++ {
		next <- true
	}

	var ps []uint64
	for i := uint64(2); i <= max; i++ {
		if nums[i] == false {
			ps = append(ps, i)
		}
	}
	return ps
}
