package service

import (
	"context"
	"errors"
	"math"

	"github.com/balinor2017/interview-order/models"
	"github.com/balinor2017/interview-order/pubsub"
	"github.com/balinor2017/interview-order/repository"
	"github.com/balinor2017/interview-order/util"
	log "github.com/sirupsen/logrus"

	"time"

	"github.com/balinor2017/interview-order/flags"
)

// OrderServiceInterface define order service interface
// include GetOrderByID, CreateOrder, GetTripByOrderID, UpdateOrderDispacherStatus
type OrderServiceInterface interface {
	GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error)
	CreateOrder(ctx context.Context, order models.Order) (*models.Order, error)
	GetTripByOrderID(ctx context.Context, orderID int64) (*models.OrderTrip, error)
	UpdateOrderDispacherStatus(ctx context.Context, order models.OrderDispacherStatus) error
}

// OrderService Implement service with empty struct
type OrderService struct {
}

// UpdateOrderDispacherStatus ...
func (OrderService) UpdateOrderDispacherStatus(ctx context.Context, order models.OrderDispacherStatus) error {
	if order.OrderID == 0 {
		log.Info("you request data not exit")
		return nil
	}
	t := util.ToTimeString(time.Now())
	repo := repository.GetOrderRepository()
	// update rs_order status
	err := repo.UpdateOrderDispacherStatus(order)
	if err != nil {
		log.Error(err)
		return err
	}
	// get order detail by order ID
	or, _ := repo.GetOrderByID(order.OrderID)

	notification := models.UserNotification{}
	notification.UserID = or.CustomerID
	switch order.Status {
	case flags.Booked:
		notification.Message = "Your order has been booked"
		err = publisher.Publish(flags.SendingNotification, notification)
		if err != nil {
			log.Errorf("Set topic error: %s", err.Error())
		}

		orderTrip := models.OrderTrip{OrderID: or.OrderID, CustomerID: or.CustomerID, DriverID: order.DriverID, StartTime: t, PickupLatitude: or.PickupLatitude, PickupLongitude: or.PickupLongitude, Status: order.Status, CreatedAt: t}
		err := repo.WriteOrderTrip(orderTrip)
		if err != nil {
			log.Error(err)
			return err
		}
	case flags.Arrived:
		notification.Message = "Arrived"
		err = publisher.Publish(flags.SendingNotification, notification)
		if err != nil {
			log.Errorf("Set topic error: %s", err.Error())
		}
	case flags.StartTrip:
		notification.Message = "You will start the trip"
		err = publisher.Publish(flags.SendingNotification, notification)
		if err != nil {
			log.Errorf("Set topic error: %s", err.Error())
		}

		vacantDistance := earthDistance(order.DriverLongitude, order.DriverLatitude, or.PickupLongitude, or.PickupLatitude)
		orderTrip := models.OrderTrip{StartTime: t, VacantDistance: vacantDistance, PickupLatitude: order.DriverLatitude, PickupLongitude: order.DriverLongitude, Status: order.Status, UpdatedAt: t, OrderID: order.OrderID}
		err := repo.ModifyOrderTrip(orderTrip)
		if err != nil {
			log.Error(err)
			return err
		}
	case flags.EndTrip:
		notification.Message = "Your trip has ended"
		err = publisher.Publish(flags.SendingNotification, notification)
		if err != nil {
			log.Errorf("Set topic error: %s", err.Error())
		}

		orderTrip, _ := repo.GetOrderTripByOrderID(order.OrderID)
		engagedDistance := earthDistance(order.DriverLongitude, order.DriverLatitude, orderTrip.PickupLongitude, orderTrip.PickupLatitude)
		tripInfo := models.OrderTrip{EndTime: t, EngagedDistance: engagedDistance, DropoffLatitude: order.DriverLatitude, DropoffLongitude: order.DriverLongitude, Status: order.Status, UpdatedAt: t, OrderID: order.OrderID}
		err := repo.ModifyOrderTrip(tripInfo)
		if err != nil {
			log.Error(err)
			return err
		}
	default:
		errors.New("order status is invalid")
	}

	return nil

}

// GetTripByOrderID ...
func (OrderService) GetTripByOrderID(ctx context.Context, orderID int64) (*models.OrderTrip, error) {
	repo := repository.GetOrderRepository()
	order, err := repo.GetOrderTripByOrderID(orderID)
	if err != nil {
		log.Error(err)
	}

	return order, err
}

// GetOrderByID ...
func (OrderService) GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error) {
	repo := repository.GetOrderRepository()
	order, err := repo.GetOrderByID(orderID)

	if order.CustomerID != 0 {
		// Get Customer Info From Customer Service
	}

	if order.DriverID != 0 {
		// Get Driver Info From Driver Service
	}
	if err != nil {
		log.Error(err)
	}

	return order, err
}

// CreateOrder ...
func (OrderService) CreateOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	repo := repository.GetOrderRepository()
	orderInfo, err := repo.WriteOrder(order)
	if err != nil {
		log.Error(err)
	}

	return orderInfo, err
}

func handleSendingNotification(data []byte, _ ...interface{}) {
	var notification models.UserNotification
	err := pubsub.DecodeTextPlain(data, &notification)
	if err != nil {
		log.Error(err)
	}

	log.Debug(notification)

	//TODO: get user device token and push notification to notification platform with user info.
}

func earthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := float64(6371000)
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}
