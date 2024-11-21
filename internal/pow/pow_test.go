package pow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsHashCorrect(t *testing.T) {
	type args struct {
		hash       string
		zerosCount int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test positive",
			args: args{
				hash:       "lmlmllnnl",
				zerosCount: 20,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, nonce, err := ComputeHashcash(GenerateChallenge(32), 1000000, 20)
			if err != nil {
				t.Errorf("ComputeHashcash() error = %v", err)
			}
			fmt.Println(hash, nonce)
			if got := IsHashCorrect(tt.args.hash, tt.args.zerosCount); got != tt.want {
				t.Errorf("IsHashCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
