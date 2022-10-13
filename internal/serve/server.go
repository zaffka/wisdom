package serve

import (
	"bytes"
	"context"
	"encoding/hex"
	"net"

	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

type Server struct {
	powComplexity int64
	listener      net.Listener
	log           *zap.Logger
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
		if err != nil {
			s.log.Error("failed to accept a connection", zap.Error(err))

			return
		}

		go s.handleCall(ctx, conn)
	}
}

func (s *Server) handleCall(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	log := s.log.With(zap.String("client_address", conn.RemoteAddr().String()))
	log.Debug("handling a call")

	block, err := pow.NewBlock(s.powComplexity)
	if err != nil {
		log.Error("failed to create a pow puzzle", zap.Error(err))

		return
	}

	salt, nonce, err := block.Parse()
	if err != nil {
		log.Error("failed to parse a block", zap.Error(err))

		return
	}

	// Sending a puzzle as a hex-encoded string.
	hashedBlock := block.Sha256()
	n, err := conn.Write(bytes.Join([][]byte{salt, hashedBlock}, []byte{}))
	if err != nil {
		log.Error("failed to write a puzzle to remote", zap.Error(err))

		return
	}
	log.Debug("puzzle sent",
		zap.Int("bytes_written", n),
		zap.String("hashed_block_hex", hex.EncodeToString(hashedBlock)),
		zap.String("salt_hex", hex.EncodeToString(salt)),
		zap.Uint64("nonce", nonce),
	)
}
