package debug

import (
	"context"
	"fmt"
	"net"

	"github.com/d2verb/gemu/pkg/debug/pb"
	"github.com/d2verb/gemu/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DebugServer struct {
	port      int
	ch        chan any
	debugMode bool

	pb.UnimplementedHealthCheckerServer
	pb.UnimplementedDebuggerServer
}

func NewDebugServer(port int, ch chan any, debugMode bool) *DebugServer {
	return &DebugServer{
		port:      port,
		ch:        ch,
		debugMode: debugMode,
	}
}

func (d *DebugServer) Hi(cxt context.Context, req *pb.HiRequest) (*pb.HiReply, error) {
	return &pb.HiReply{}, nil
}

func (d *DebugServer) Next(cxt context.Context, req *pb.NextRequest) (*pb.NextReply, error) {
	d.ch <- req
	<-d.ch
	return &pb.NextReply{}, nil
}

func (d *DebugServer) Start(ctx context.Context, cancel context.CancelFunc) {
	if !d.debugMode {
		return
	}

	log.Debugf("Starting debug server (port: %d)...\n", d.port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", d.port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterHealthCheckerServer(s, d)
	pb.RegisterDebuggerServer(s, d)

	reflection.Register(s)

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
