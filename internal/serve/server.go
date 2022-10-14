package serve

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/zaffka/wisdom/assets"
	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

const sendingDeadline = 5 * time.Second

type Server struct {
	powComplexity  int64
	makePowBlockFn func(int64) (pow.Block, error)
	listener       net.Listener
	log            *zap.Logger
}

func NewServer(optFns ...OptFn) *Server {
	s := &Server{}

	for _, fn := range optFns {
		fn(s)
	}

	return s
}

func (s *Server) Run(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			s.log.Error("failed to accept a connection, halting the app", zap.Error(err))

			return
		}

		go func(conn ConnHandler) {
			log := s.log.With(zap.String("client_address", conn.RemoteAddr().String()))
			log.Debug("handling a call")

			quote, err := s.handleCall(ctx, conn)
			if err != nil {
				log.Error("failed to handle call", zap.Error(err))
			} else {
				log.Debug("puzzle successfully solved, quote sent", zap.String("quote", quote))
			}

			if err := conn.Close(); err != nil {
				log.Error("failed to call Close on a connection", zap.Error(err))
			}
		}(conn)
	}
}

func (s *Server) handleCall(ctx context.Context, conn ConnHandler) (string, error) {
	block, err := s.makePowBlockFn(s.powComplexity)
	if err != nil {
		return "", fmt.Errorf("failed to create a pow puzzle: %w", err)
	}

	salt, _, err := block.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to parse a block: %w", err)
	}

	// Sending a puzzle.
	// It consist of a salt (as is) and a hashed salt+nonce block.
	hashedBlock := block.Sha256()
	_, err = conn.Write(bytes.Join([][]byte{salt, hashedBlock}, []byte{}))
	if err != nil {
		return "", fmt.Errorf("failed to write a puzzle to remote: %w", err)
	}

	// Reading the solution.
	// The client must find the nonce using salt and a hash of salt+nonce block.
	if err := conn.SetReadDeadline(time.Now().Add(sendingDeadline)); err != nil {
		return "", fmt.Errorf("failed to set a deadline for puzzle solving: %w", err)
	}

	nonceResp := make([]byte, 8) // expected BigEndian bytes order within
	_, err = conn.Read(nonceResp)
	if err != nil {
		return "", fmt.Errorf("failed to read a pow result: %w", err)
	}

	// Comparing with a nonce we have.
	if !bytes.Equal(nonceResp, block[8:]) {
		return "", fmt.Errorf("got wrong pow result: %w", err)
	}

	// Sending a random quote back.
	quote := assets.RandomQuote()
	_, err = conn.Write([]byte(quote))
	if err != nil {
		return "", fmt.Errorf("failed to write a wisdom quote to remote: %w", err)
	}

	return quote, nil
}
