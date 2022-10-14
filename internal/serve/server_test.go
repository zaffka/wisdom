package serve

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zaffka/wisdom/internal/serve/mocks"
	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

var errTest = errors.New("testerr")

func TestServer_handleCall(t *testing.T) {
	t.Parallel()

	log := zap.NewNop()
	block, err := pow.NewBlock(100)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = block.Parse()
	if err != nil {
		t.Fatal(err)
	}

	dummyBlockFn := func(int64) (pow.Block, error) {
		return block, nil
	}

	type fields struct {
		powComplexity  int64
		makePowBlockFn func(int64) (pow.Block, error)
		log            *zap.Logger
	}
	type args struct {
		connHandlerFn func(*testing.T, *mocks.ConnHandler)
		expectedErr   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				powComplexity:  100000,
				makePowBlockFn: dummyBlockFn,
				log:            log,
			},
			args: args{
				connHandlerFn: func(t *testing.T, m *mocks.ConnHandler) {
					m.On("Write", mock.Anything).Return(8, nil).Once()
					m.On("SetReadDeadline", mock.Anything).Return(nil)
					m.On("Read", mock.Anything).Return(8, nil).Run(func(args mock.Arguments) {
						arg := args.Get(0).([]byte)
						copy(arg, block[8:])
					})
					m.On("Write", mock.Anything).Return(32, nil).Once()
				},
				expectedErr: nil,
			},
		},
		{
			name: "fail on sending quote",
			fields: fields{
				powComplexity:  100000,
				makePowBlockFn: dummyBlockFn,
				log:            log,
			},
			args: args{
				connHandlerFn: func(t *testing.T, m *mocks.ConnHandler) {
					m.On("Write", mock.Anything).Return(8, nil).Once()
					m.On("SetReadDeadline", mock.Anything).Return(nil)
					m.On("Read", mock.Anything).Return(8, nil).Run(func(args mock.Arguments) {
						arg := args.Get(0).([]byte)
						copy(arg, block[8:])
					})
					m.On("Write", mock.Anything).Return(0, errTest).Once()
				},
				expectedErr: fmt.Errorf("failed to write a wisdom quote to remote: %w", errTest),
			},
		},
		{
			name: "fail on reading pow result",
			fields: fields{
				powComplexity:  100000,
				makePowBlockFn: dummyBlockFn,
				log:            log,
			},
			args: args{
				connHandlerFn: func(t *testing.T, m *mocks.ConnHandler) {
					m.On("Write", mock.Anything).Return(8, nil).Once()
					m.On("SetReadDeadline", mock.Anything).Return(nil)
					m.On("Read", mock.Anything).Return(0, errTest)
				},
				expectedErr: fmt.Errorf("failed to read a pow result: %w", errTest),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Server{
				powComplexity:  tt.fields.powComplexity,
				makePowBlockFn: tt.fields.makePowBlockFn,
				log:            tt.fields.log,
			}

			connHandlerMock := mocks.NewConnHandler(t)
			tt.args.connHandlerFn(t, connHandlerMock)

			quote, err := s.handleCall(context.Background(), connHandlerMock)
			if err != nil {
				require.Equal(t, tt.args.expectedErr, err)
			} else {
				require.NotEmpty(t, quote)
			}

		})
	}
}
