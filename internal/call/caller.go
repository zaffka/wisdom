package call

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

const endOfLine = '\n'

type Caller struct {
	protocol string
	addr     string
	log      *zap.Logger
}

func NewCaller(optFns ...OptFn) *Caller {
	c := &Caller{}

	for _, fn := range optFns {
		fn(c)
	}

	return c
}

func (c *Caller) Run(ctx context.Context) error {
	startTime := time.Now()

	conn, err := net.Dial(c.protocol, c.addr)
	if err != nil {
		return fmt.Errorf("failed to call the server: %w", err)
	}

	puzzleBuf := make([]byte, 40)
	read, err := conn.Read(puzzleBuf)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	c.log.Debug("data read", zap.Int("bytes_read", read))

	nonce, err := pow.FindNonce(puzzleBuf[8:], puzzleBuf[:8])
	if err != nil {
		return fmt.Errorf("failed to find a nonce: %w", err)
	}

	c.log.Info("puzzle solved", zap.Uint64("nonce", nonce), zap.String("time_spent", time.Now().Sub(startTime).String()))

	return nil
}
