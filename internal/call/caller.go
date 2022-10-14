package call

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

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
		return fmt.Errorf("failed to dial the server: %w", err)
	}

	// Reading puzzle from the server.
	puzzleBuf := make([]byte, 40)
	_, err = conn.Read(puzzleBuf)
	if err != nil {
		return fmt.Errorf("failed to read data from the server: %w", err)
	}

	// Solving the puzzle.
	nonce, err := pow.FindNonce(puzzleBuf[8:], puzzleBuf[:8])
	if err != nil {
		return fmt.Errorf("failed to find a nonce: %w", err)
	}

	c.log.Info("puzzle solved", zap.Uint64("nonce", nonce), zap.String("time_spent", time.Since(startTime).String()))

	// Sending nonce back.
	nonceResp := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceResp, nonce)
	_, err = conn.Write(nonceResp)
	if err != nil {
		return fmt.Errorf("failed to write a nonce to the server: %w", err)
	}

	// Reading wisdom quote.
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return fmt.Errorf("failed to set a reading deadline: %w", err)
	}

	quoteBuf := &bytes.Buffer{}
	_, err = io.Copy(quoteBuf, conn)
	if err != nil {
		return fmt.Errorf("failed to read a quote bytes from the server: %w", err)
	}

	c.log.Info("got wisdom quote", zap.ByteString("data", quoteBuf.Bytes()))

	return nil
}
