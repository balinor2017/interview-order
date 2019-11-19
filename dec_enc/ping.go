package dec_enc

import (
	"context"

	"net/http"

	jsonpb "github.com/golang/protobuf/jsonpb"
	empty "github.com/golang/protobuf/ptypes/empty"

	"github.com/balinor2017/interview-order/pb"
)

func DecodeGRPCPingRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &empty.Empty{}, nil
}

func EncodeGRPCPingResponse(_ context.Context, resp interface{}) (interface{}, error) {
	return resp.(*pb.Response), nil
}

func DecodeHTTPPingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &empty.Empty{}, nil
}

func EncodeHTTPPingResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	ma := jsonpb.Marshaler{}
	return ma.Marshal(w, resp.(*pb.Response))
}
