package gameboy

import (
	"context"
	"fmt"
	"net"

	"github.com/d2verb/gemu/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/d2verb/gemu/pkg/gameboy/debug"
)

type DebugServer struct {
	port      int
	debugMode bool
	pb.UnimplementedHealthCheckerServer
}

func newDebugServer(port int, debugMode bool) *DebugServer {
	return &DebugServer{
		port:      port,
		debugMode: debugMode,
	}
}

func (d *DebugServer) Hi(cxt context.Context, req *pb.HiRequest) (*pb.HiReply, error) {
	return &pb.HiReply{}, nil
}

func (d *DebugServer) start(ctx context.Context, cancel context.CancelFunc) {
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
	reflection.Register(s)

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
