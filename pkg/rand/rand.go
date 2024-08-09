package rand

import (
	"math/rand/v2"
)

// PermFrom1 generates random permutation from 1 to n
func PermFrom1(n int) []int {
	perm := rand.Perm(n)
	for i := 0; i < len(perm); i++ {
		perm[i]++
	}
	return perm
}
