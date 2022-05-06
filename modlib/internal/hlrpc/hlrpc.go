package hlrpc

import (
	"context"
	"hlinspect/internal/handlers"
	"hlinspect/internal/hlrpc/schema"
	"hlinspect/internal/logs"
	"io"
	"net"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
)

func Serve(handler *handlers.HLRPCHandler) error {
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

func serveHLRPC(rwc io.ReadWriteCloser, handler *handlers.HLRPCHandler) {
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
	handler *handlers.HLRPCHandler
}

func (s *halflifeServer) GetFullPlayerState(ctx context.Context, call schema.HalfLife_getFullPlayerState) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}

	fullState, err := schema.NewRootFullPlayerState(seg)
	if err != nil {
		return err
	}

	err = s.handler.GetFullPlayerState(ctx, &fullState)
	if err != nil {
		return err
	}

	return res.SetState(fullState)
}

func (s *halflifeServer) StartInputControl(ctx context.Context, call schema.HalfLife_startInputControl) error {
	return s.handler.StartInputControl(ctx)
}

func (s *halflifeServer) StopInputControl(ctx context.Context, call schema.HalfLife_stopInputControl) error {
	return s.handler.StopInputControl(ctx)
}

func (s *halflifeServer) InputStep(ctx context.Context, call schema.HalfLife_inputStep) error {
	cmd, err := call.Args().Cmd()
	if err != nil {
		return err
	}

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}

	fullState, err := schema.NewRootFullPlayerState(seg)
	if err != nil {
		return err
	}

	err = s.handler.InputStep(ctx, &cmd, &fullState)
	if err != nil {
		return err
	}

	return res.SetState(fullState)
}
