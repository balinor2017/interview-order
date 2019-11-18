package dec_enc

import (
	"context"

	"net/http"

	"github.com/golang/protobuf/jsonpb"

	"github.com/interview-order/pb"
)

// UpdateOrderDispacherStatus
func DecodeHTTPUpdateOrderDispacherStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.OrderDispacherStatusRequest
	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeHTTPUpdateOrderDispacherStatusResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	ma := jsonpb.Marshaler{}
	return ma.Marshal(w, resp.(*pb.Response))
}

// Get trip by order ID
func DecodeHTTPGetTripByOrderIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetTripByOrderIDRequest
	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeHTTPGetTripByOrderIDResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	ma := jsonpb.Marshaler{}
	return ma.Marshal(w, resp.(*pb.Response))
}

// Get order By ID
func DecodeHTTPGetOrderByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetOrderByIDRequest
	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeHTTPGetOrderByIDResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	ma := jsonpb.Marshaler{}
	return ma.Marshal(w, resp.(*pb.Response))
}

// Create order
func DecodeHTTPCreateOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.Order
	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeHTTPCreateOrderResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	ma := jsonpb.Marshaler{}
	return ma.Marshal(w, resp.(*pb.Response))
}
