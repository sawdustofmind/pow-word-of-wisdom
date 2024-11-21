package pow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkComputeHashcash(b *testing.B) {
	maxIterations := 100000000

	for _, l := range []int{32, 64} {
		for _, zerosCount := range []int{3, 4, 5, 6} {
			b.Run(fmt.Sprintf("l = %d, zeros = %d", l, zerosCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _, err := ComputeHashcash(GenerateChallenge(l), maxIterations, zerosCount)
					require.NoError(b, err)
				}
			})
		}
	}

}

func TestGenerateChallenge(t *testing.T) {
	result := GenerateChallenge(15)
	assert.NotEmpty(t, result)
	assert.Len(t, result, 30)
}
