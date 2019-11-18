package server

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/interview-order/dec_enc"
	"github.com/interview-order/endpoint"
	"github.com/interview-order/pb"
)

type grpcServer struct {
	ping grpctransport.Handler
}

// implement LoremServer Interface in lorem.pb.go
func (s *grpcServer) Ping(ctx context.Context, req *empty.Empty) (*pb.Response, error) {
	_, resp, err := s.ping.ServeGRPC(ctx, nil)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Response), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint endpoint.Endpoints) pb.PingServer {
	return &grpcServer{
		ping: grpctransport.NewServer(
			endpoint.PingEndpoint,
			dec_enc.DecodeGRPCPingRequest,
			dec_enc.EncodeGRPCPingResponse,
		),
	}
}
