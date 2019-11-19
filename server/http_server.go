package server

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/balinor2017/interview-order/dec_enc"
	"github.com/balinor2017/interview-order/endpoint"
)

func NewHttpServer(ctx context.Context, endpoints endpoint.Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/ping",
		httptransport.NewServer(
			endpoints.PingEndpoint,
			dec_enc.DecodeHTTPPingRequest,
			dec_enc.EncodeHTTPPingResponse,
		))

	m.Handle("/order",
		httptransport.NewServer(
			endpoints.GetOrderByIDEndpoint,
			dec_enc.DecodeHTTPGetOrderByIDRequest,
			dec_enc.EncodeHTTPGetOrderByIDResponse,
		))

	m.Handle("/order/create",
		httptransport.NewServer(
			endpoints.CreateOrderEndpoint,
			dec_enc.DecodeHTTPCreateOrderRequest,
			dec_enc.EncodeHTTPCreateOrderResponse,
		))

	m.Handle("/order/gettripbyorderid",
		httptransport.NewServer(
			endpoints.GetTripByOrderIDEndpoint,
			dec_enc.DecodeHTTPGetTripByOrderIDRequest,
			dec_enc.EncodeHTTPGetTripByOrderIDResponse,
		))

	m.Handle("/order/dispacherstatus",
		httptransport.NewServer(
			endpoints.UpdateOrderDispacherStatusEndpoint,
			dec_enc.DecodeHTTPUpdateOrderDispacherStatusRequest,
			dec_enc.EncodeHTTPUpdateOrderDispacherStatusResponse,
		))

	return m
}
