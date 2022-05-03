package hlrpc

import (
	"context"
	"hlinspect/internal/hlrpc/schema"
	"hlinspect/internal/logs"
	"io"
	"net"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
)

type Handler interface {
	GetFullPlayerState(ctx context.Context, resp *schema.FullPlayerState) error
}

func Serve(handler Handler) error {
	// FIXME: make this configurable instead of accepting from all
	// nolint:gosec
	ln, err := net.Listen("tcp", "0.0.0.0:32002")
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		logs.DLLLog.Infof("hlrpc received a new connection: %v", conn.RemoteAddr())

		go serveHLRPC(conn, handler)
	}
}

func serveHLRPC(rwc io.ReadWriteCloser, handler Handler) {
	defer rwc.Close()

	main := schema.HalfLife_ServerToClient(&halflifeServer{
		handler: handler,
	}, nil)

	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		BootstrapClient: main.Client,
	})
	defer conn.Close()

	<-conn.Done()
}

type halflifeServer struct {
	handler Handler
}

func (s *halflifeServer) GetFullPlayerState(ctx context.Context, call schema.HalfLife_getFullPlayerState) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}

	fullState, err := schema.NewRootFullPlayerState(seg)
	if err != nil {
		panic(err)
	}

	err = s.handler.GetFullPlayerState(ctx, &fullState)
	if err != nil {
		return err
	}

	return res.SetState(fullState)
}
