package serve

import (
	"net"

	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

type OptFn func(*Server)

func WithListener(listener net.Listener) OptFn {
	return func(s *Server) {
		s.listener = listener
	}
}

func WithLogger(logger *zap.Logger) OptFn {
	return func(s *Server) {
		s.log = logger
	}
}

func WithInitialPoWComplexity(cmpx int64) OptFn {
	return func(s *Server) {
		s.powComplexity = cmpx
	}
}

func WithPoWBlockFunc(powFn func(int64) (pow.Block, error)) OptFn {
	return func(s *Server) {
		s.makePowBlockFn = powFn
	}
}
