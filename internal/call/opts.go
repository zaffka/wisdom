package call

import "go.uber.org/zap"

type OptFn func(*Caller)

func WithServerAddr(address string) OptFn {
	return func(c *Caller) {
		c.addr = address
	}
}

func WithLogger(logger *zap.Logger) OptFn {
	return func(c *Caller) {
		c.log = logger
	}
}

func WithProtocol(protocol string) OptFn {
	return func(c *Caller) {
		c.protocol = protocol
	}
}
