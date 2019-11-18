package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang/protobuf/ptypes"
	"github.com/interview-order/models"
	"github.com/interview-order/pb"
	"github.com/interview-order/service"
)

// Wrapping Endpoints as a Service implementation.
// UpdateOrderDispacherStatus
func (e Endpoints) UpdateOrderDispacherStatus(ctx context.Context, req *pb.OrderDispacherStatusRequest) (*pb.GetOrderByIDResponse, error) {
	resp, err := e.CreateOrderEndpoint(ctx, req)
	if err != nil {
		return &pb.GetOrderByIDResponse{}, err
	}
	res := resp.(*pb.GetOrderByIDResponse)
	return res, nil
}

func MakeUpdateOrderDispacherStatusEndpoint(svc service.OrderServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.OrderDispacherStatusRequest)
		// request
		orderReq := models.OrderDispacherStatus{
			OrderID:         req.OrderId,
			DriverID:        req.DriverId,
			DriverLongitude: req.DriverLongitude,
			DriverLatitude:  req.DriverLatitude,
			Status:          req.Status,
		}

		err := svc.UpdateOrderDispacherStatus(ctx, orderReq)

		data, err := ptypes.MarshalAny(&pb.GeneraResponse{
			Msg: "update success",
		})

		code := 200
		msg := "ok"
		if err != nil {
			code = 500
			msg = err.Error()
		}

		return &pb.Response{
			Code:     int32(code),
			Message:  msg,
			Response: data,
		}, nil
	}
}

// GetTripByOrderID
func (e Endpoints) GetTripByOrderID(ctx context.Context, req *pb.GetTripByOrderIDRequest) (*pb.GetTripByOrderIDResponse, error) {
	resp, err := e.GetTripByOrderIDEndpoint(ctx, req)
	if err != nil {
		return &pb.GetTripByOrderIDResponse{}, err
	}
	res := resp.(*pb.GetTripByOrderIDResponse)
	return res, nil
}

func MakeGetTripByOrderIDEndpoint(svc service.OrderServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetTripByOrderIDRequest)
		res, err := svc.GetTripByOrderID(ctx, req.Id)
		orderTrip := &pb.OrderTrip{
			TripId:           res.TripID,
			OrderId:          res.OrderID,
			CustomerId:       res.CustomerID,
			DriverId:         res.DriverID,
			StartTime:        res.StartTime,
			EndTime:          res.EndTime,
			VacantDistance:   res.VacantDistance,
			EngagedDistance:  res.EngagedDistance,
			PickupLatitude:   res.PickupLatitude,
			PickupLongitude:  res.PickupLongitude,
			DropoffLatitude:  res.DropoffLatitude,
			DropoffLongitude: res.DropoffLongitude,
			Status:           res.Status,
			CreatedAt:        res.CreatedAt,
			UpdatedAt:        res.UpdatedAt,
		}

		data, err := ptypes.MarshalAny(&pb.GetTripByOrderIDResponse{
			Data: orderTrip,
		})

		code := 200
		msg := "ok"
		if err != nil {
			code = 500
			msg = err.Error()
		}

		return &pb.Response{
			Code:     int32(code),
			Message:  msg,
			Response: data,
		}, nil
	}
}

// get order By ID
func (e Endpoints) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	resp, err := e.GetOrderByIDEndpoint(ctx, req)
	if err != nil {
		return &pb.GetOrderByIDResponse{}, err
	}
	res := resp.(*pb.GetOrderByIDResponse)
	return res, nil
}

func MakeGetOrderByIDEndpoint(svc service.OrderServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetOrderByIDRequest)

		res, err := svc.GetOrderByID(ctx, req.Id)
		order := &pb.Order{
			OrderId:          res.OrderID,
			CustomerId:       res.CustomerID,
			CustomerName:     res.CustomerName,
			CustomerMobile:   res.CustomerMobile,
			DriverId:         res.DriverID,
			DriverName:       res.DriverName,
			DriverMobile:     res.DriverMobile,
			VehicleId:        res.VehicleID,
			VehicleNo:        res.VehicleNo,
			PickTime:         res.PickupTime,
			PickupLatitude:   res.PickupLatitude,
			PickupLongitude:  res.PickupLongitude,
			DropoffLatitude:  res.DropoffLatitude,
			DropoffLongitude: res.DropoffLongitude,
			Status:           res.Status,
			CreatedAt:        res.CreatedAt,
			UpdatedAt:        res.UpdatedAt,
		}

		data, err := ptypes.MarshalAny(&pb.GetOrderByIDResponse{
			Data: order,
		})

		code := 200
		msg := "ok"
		if err != nil {
			code = 500
			msg = err.Error()
		}

		return &pb.Response{
			Code:     int32(code),
			Message:  msg,
			Response: data,
		}, nil
	}
}

// create order

func (e Endpoints) GreateOrder(ctx context.Context, req *pb.Order) (*pb.GetOrderByIDResponse, error) {
	resp, err := e.CreateOrderEndpoint(ctx, req)
	if err != nil {
		return &pb.GetOrderByIDResponse{}, err
	}
	res := resp.(*pb.GetOrderByIDResponse)
	return res, nil
}

func MakeCreateOrderEndpoint(svc service.OrderServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.Order)
		// request
		orderReq := models.Order{
			OrderID:          req.OrderId,
			CustomerID:       req.CustomerId,
			CustomerName:     req.CustomerName,
			CustomerMobile:   req.CustomerMobile,
			DriverID:         req.DriverId,
			DriverName:       req.DriverName,
			DriverMobile:     req.DriverMobile,
			VehicleID:        req.VehicleId,
			VehicleNo:        req.VehicleNo,
			PickupTime:       req.PickTime,
			PickupLatitude:   req.PickupLatitude,
			PickupLongitude:  req.PickupLongitude,
			DropoffLatitude:  req.DropoffLatitude,
			DropoffLongitude: req.DropoffLongitude,
			Status:           req.Status,
			CreatedAt:        req.CreatedAt,
			UpdatedAt:        req.UpdatedAt,
		}

		res, err := svc.CreateOrder(ctx, orderReq)

		// response
		order := &pb.Order{
			OrderId:          res.OrderID,
			CustomerId:       res.CustomerID,
			CustomerName:     res.CustomerName,
			CustomerMobile:   res.CustomerMobile,
			DriverId:         res.DriverID,
			DriverName:       res.DriverName,
			DriverMobile:     res.DriverMobile,
			VehicleId:        res.VehicleID,
			VehicleNo:        res.VehicleNo,
			PickTime:         res.PickupTime,
			PickupLatitude:   res.PickupLatitude,
			PickupLongitude:  res.PickupLongitude,
			DropoffLatitude:  res.DropoffLatitude,
			DropoffLongitude: res.DropoffLongitude,
			Status:           res.Status,
			CreatedAt:        res.CreatedAt,
			UpdatedAt:        res.UpdatedAt,
		}

		data, err := ptypes.MarshalAny(&pb.GetOrderByIDResponse{
			Data: order,
		})

		code := 200
		msg := "ok"
		if err != nil {
			code = 500
			msg = err.Error()
		}

		return &pb.Response{
			Code:     int32(code),
			Message:  msg,
			Response: data,
		}, nil
	}
}
