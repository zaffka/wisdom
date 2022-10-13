package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
)

const (
	BlockSize = 16
)

var (
	ErrWrongBlockSize        = errors.New("wrong size of the Block")
	ErrNonceNotFound         = errors.New("failed to find a Nonce")
	ErrIntUintRepresentation = errors.New("failed to represent maxInt as unsigned uint64")
)

// Block is a slice of bytes of 16-bytes size (BlockSize).
// First 8 bytes of a Block holds a random-generated salt.
// We need this to prevent Rainbow table attack https://en.wikipedia.org/wiki/Rainbow_table.
// Second 8 bytes of a Block holds a Nonce - random-generated uint64 value.
// To solve a Proof-of-Work puzzle you need to find that Nonce value using FindNonce function.
type Block []byte

// CheckSize checks if the Block has a correct size of 16 bytes.
func (b Block) CheckSize() error {
	if len(b) != BlockSize {
		return fmt.Errorf("%w, got %d size, need %d", ErrWrongBlockSize, len(b), BlockSize)
	}

	return nil
}

// Parse extracts raw salt and nonce from a Block.
func (b Block) Parse() ([]byte, uint64, error) {
	if err := b.CheckSize(); err != nil {
		return nil, 0, err
	}

	return b[:8], binary.BigEndian.Uint64(b[8:]), nil
}

// Sha256 hashes the entire Block to create a puzzle to be solved.
func (b Block) Sha256() []byte {
	hashedBlock := sha256.Sum256(b)

	return hashedBlock[:]
}

// NewBlock creates a new Block with 8-bytes sized salt and 8-bytes sized nonce.
// It uses cryptographically secure bytes generator and maxInt value to limit Nonce max value.
func NewBlock(maxInt int64) (Block, error) {
	block := make(Block, BlockSize)
	_, err := io.ReadFull(rand.Reader, block[:8])
	if err != nil {
		return nil, fmt.Errorf("failed to create a random bytes: %w", err)
	}

	bi := big.NewInt(maxInt)
	if !bi.IsUint64() || bi.Sign() <= 0 {
		return nil, fmt.Errorf("%w, got int64(%s)", ErrIntUintRepresentation, bi.String())
	}

	i, err := rand.Int(rand.Reader, bi)
	if err != nil {
		return nil, fmt.Errorf("failed to create a random int: %w", err)
	}

	binary.BigEndian.PutUint64(block[8:], i.Uint64())

	return block, nil
}

// IsValid makes hash of a Block and compares it with the provided one.
func IsValid(block Block, comparableHash []byte) bool {
	if err := block.CheckSize(); err != nil {
		return false
	}

	return bytes.Equal(comparableHash, block.Sha256())
}

// FindNonce is picking up a Nonce, hashing its value with already known salt.
// When the resulting hash is the same as provided - puzzle solved.
func FindNonce(hash, salt []byte) (uint64, error) {
	nonceBlock := make([]byte, 8)
	for i := int64(0); i < math.MaxInt64; i++ {
		binary.BigEndian.PutUint64(nonceBlock, uint64(i))
		possibleBlock := Block(bytes.Join([][]byte{salt, nonceBlock}, []byte{}))
		if bytes.Equal(hash, possibleBlock.Sha256()) {
			return binary.BigEndian.Uint64(nonceBlock), nil
		}
	}

	return 0, ErrNonceNotFound
}
