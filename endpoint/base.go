package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	empty "github.com/golang/protobuf/ptypes/empty"

	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/interview-order/pb"
	"github.com/interview-order/service"
)

// endpoints wrapper
type Endpoints struct {
	PingEndpoint                       endpoint.Endpoint
	GetOrderByIDEndpoint               endpoint.Endpoint
	CreateOrderEndpoint                endpoint.Endpoint
	GetTripByOrderIDEndpoint           endpoint.Endpoint
	UpdateOrderDispacherStatusEndpoint endpoint.Endpoint
}

// Wrapping Endpoints as a Service implementation.
func (e Endpoints) Ping(ctx context.Context, req *empty.Empty) (string, error) {
	resp, err := e.PingEndpoint(ctx, nil)
	if err != nil {
		return "", err
	}
	res := resp.(*pb.PingResponse)
	return res.Text, nil
}

// creating Lorem Ipsum Endpoint
func MakePingEndpoint(svc service.IPingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*empty.Empty)

		txt, error := svc.Ping(ctx, req)

		data, error := ptypes.MarshalAny(&pb.PingResponse{
			Text: txt,
		})

		return &pb.Response{
			Code:     200,
			Message:  "ok",
			Response: data,
		}, error
	}
}
