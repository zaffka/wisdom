package pow_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zaffka/wisdom/pkg/pow"
)

func TestNewBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		maxInt int64
		block  pow.Block
		err    error
	}{
		{
			name:   "max int as negative int",
			maxInt: -1,
			block:  nil,
			err:    fmt.Errorf("%w, got int64(%d)", pow.ErrIntUintRepresentation, -1),
		},
		{
			name:   "max int as zero",
			maxInt: 0,
			block:  nil,
			err:    fmt.Errorf("%w, got int64(%d)", pow.ErrIntUintRepresentation, 0),
		},
		{
			name:   "max int as math.MaxInt64",
			maxInt: math.MaxInt64,
			block:  []byte{0},
			err:    nil,
		},
		{
			name:   "max int as normal int",
			maxInt: 5,
			block:  []byte{0},
			err:    nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			block, err := pow.NewBlock(tt.maxInt)
			if err != nil {
				require.Equal(t, err, tt.err)
			}
			if block != nil {
				require.NotNil(t, tt.block)
			}
		})
	}
}

func TestFindNonce(t *testing.T) {
	maxInt := int64(300000)
	block, err := pow.NewBlock(maxInt)
	require.NoError(t, err)

	hashedBlock := block.Sha256()
	salt, nonce, err := block.Parse()
	require.NoError(t, err)

	foundNonce, err := pow.FindNonce(hashedBlock, salt)
	require.NoError(t, err)
	require.Equal(t, nonce, foundNonce)
}
