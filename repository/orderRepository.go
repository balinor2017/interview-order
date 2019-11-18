package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/interview-order/flags"
	"github.com/interview-order/models"
	"github.com/interview-order/util"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	creatingOrderSQL = `insert into rs_order(pickup_time,customer_id ,driver_id,pickup_longitude ,pickup_latitude ,dropoff_longitude ,dropoff_latitude,status,created_at ,updated_at) values (?,?,?,?,?,?,?,?,?,?)`

	updatingOrderSQL = `update rs_order set pickup_time=?,customer_id=?,driver_id=?,pickup_longitude=? ,pickup_latitude=?,dropoff_longitude=?,dropoff_latitude=?,status=?,updated_at=? where order_id=?`

	updatingDispacherStatusSQL = `update rs_order set driver_id=?,status=?,updated_at=? where order_id=?`

	gettingOrderByIDSQL = `SELECT order_id,pickup_time,customer_id,driver_id,
		pickup_longitude,pickup_latitude,dropoff_longitude ,dropoff_latitude,status,created_at,updated_at FROM rs_order where order_id=?`
	gettingOrderByStatusSQL = `SELECT order_id,pickup_time,customer_id,driver_id,pickup_longitude,pickup_latitude,dropoff_longitude ,dropoff_latitude,status,created_at,updated_at FROM rs_order where status=?`

	creatingOrderTripSQL = `insert into rs_trip (customer_id,order_id,driver_id,start_time,end_time,vacant_distance ,engaged_distance,
	pickup_longitude ,pickup_latitude  ,dropoff_longitude  ,dropoff_latitude,Status,created_at,updated_at) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	updatingOrderTripStartSQL = `update rs_trip set start_time=?,vacant_distance=?,pickup_longitude=? ,pickup_latitude=?,status=?,updated_at=? where order_id=?`
	updatingOrderTripEndSQL   = `update rs_trip set end_time=?,engaged_distance=?,dropoff_longitude=?,dropoff_latitude=? ,status=?,updated_at=? where order_id=?`

	gettingOrderTripByOrderIDSQL = `select  trip_id,customer_id,order_id,driver_id,start_time,end_time,vacant_distance ,engaged_distance,
	pickup_longitude ,pickup_latitude  ,dropoff_longitude  ,dropoff_latitude,Status,created_at,updated_at from rs_trip where order_id=?`
)

var (
	creatingOrderStmt             *sqlx.Stmt
	updatingOrderStmt             *sqlx.Stmt
	updatingDispacherStatusStmt   *sqlx.Stmt
	gettingOrderByIDStmt          *sqlx.Stmt
	gettingOrderByStatusStmt      *sqlx.Stmt
	creatingOrderTripStmt         *sqlx.Stmt
	updatingOrderTripStartStmt    *sqlx.Stmt
	updatingOrderTripEndStmt      *sqlx.Stmt
	gettingOrderTripByOrderIDStmt *sqlx.Stmt
)

type OrderRepository interface {
	WriteOrder(models.Order) (*models.Order, error)
	UpdateOrder(models.Order) error

	UpdateOrderDispacherStatus(order models.OrderDispacherStatus) error

	GetOrderByID(int64) (*models.Order, error)
	GetOrderByStatus(int) (*models.Order, error)

	WriteOrderTrip(models.OrderTrip) error
	ModifyOrderTrip(order models.OrderTrip) error
	GetOrderTripByOrderID(orderID int64) (*models.OrderTrip, error)
}

type orderRepository struct {
	db    *DataSource
	redis redisconection
}

var orderRepo *orderRepository

func GetOrderRepository() OrderRepository {
	if orderRepo == nil {
		orderRepo = &orderRepository{db: psqlDataSource, redis: Redis}
	}
	return orderRepo

}

func orderPrepareStmt(db *DataSource) error {
	var err error

	creatingOrderStmt, err = db.Preparex(creatingOrderSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	updatingOrderStmt, err = db.Preparex(updatingOrderSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	updatingDispacherStatusStmt, err = db.Preparex(updatingDispacherStatusSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	gettingOrderByIDStmt, err = db.Preparex(gettingOrderByIDSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	gettingOrderByStatusStmt, err = db.Preparex(gettingOrderByStatusSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	creatingOrderTripStmt, err = db.Preparex(creatingOrderTripSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	updatingOrderTripStartStmt, err = db.Preparex(updatingOrderTripStartSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	updatingOrderTripEndStmt, err = db.Preparex(updatingOrderTripEndSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	gettingOrderTripByOrderIDStmt, err = db.Preparex(gettingOrderTripByOrderIDSQL)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (repo *orderRepository) GetOrderTripByOrderID(orderID int64) (*models.OrderTrip, error) {
	order := &models.OrderTrip{}

	cacheData, err := repo.redis.getBytes(flags.CacheOrderTripPrefix + strconv.FormatInt(orderID, 10))
	if err == nil {
		err := json.Unmarshal(cacheData, order)
		if err == nil {
			return order, nil
		}
	}

	var rows *sql.Rows

	defer rows.Close()

	rows, err = gettingOrderTripByOrderIDStmt.Query(orderID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&order.TripID, &order.CustomerID, &order.OrderID, &order.DriverID,
			&order.StartTime, &order.EndTime, &order.VacantDistance, &order.EngagedDistance,
			&order.PickupLongitude, &order.PickupLatitude, &order.DropoffLongitude, &order.DropoffLatitude, &order.Status, &order.CreatedAt, &order.UpdatedAt)

		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	data, err := json.Marshal(order)
	if err == nil {
		err := repo.redis.setBytes(flags.CacheOrderTripPrefix+strconv.FormatInt(orderID, 10), []byte(data), flags.CACHETTL)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error(err)
	}

	return order, nil
}

func (repo *orderRepository) ModifyOrderTrip(orderTrip models.OrderTrip) error {

	switch orderTrip.Status {
	case flags.StartTrip:
		_, err := updatingOrderTripStartStmt.Exec(orderTrip.StartTime, orderTrip.VacantDistance, orderTrip.PickupLongitude, orderTrip.PickupLatitude, orderTrip.Status, orderTrip.UpdatedAt, orderTrip.OrderID)
		if err != nil {
			log.Error(err)
			return err
		}

	case flags.EndTrip:
		_, err := updatingOrderTripEndStmt.Exec(orderTrip.EndTime, orderTrip.EngagedDistance, orderTrip.DropoffLongitude, orderTrip.DropoffLatitude, orderTrip.Status, orderTrip.UpdatedAt, orderTrip.OrderID)
		if err != nil {
			log.Error(err)
			return err
		}

	default:
		return errors.New("orderTrip status is invalid")
	}

	repo.redis.delCache(fmt.Sprintf(flags.CacheOrderTripPrefix, orderTrip.OrderID))

	return nil
}

func (repo *orderRepository) WriteOrderTrip(orderTrip models.OrderTrip) error {

	t := util.ToTimeString(time.Now())
	result, err := creatingOrderStmt.Exec(orderTrip.CustomerID, orderTrip.OrderID, orderTrip.DriverID, orderTrip.StartTime, t,
		orderTrip.VacantDistance, orderTrip.EngagedDistance, orderTrip.PickupLongitude, orderTrip.PickupLatitude, orderTrip.DropoffLongitude, orderTrip.DropoffLatitude, orderTrip.Status, orderTrip.CreatedAt, t)

	if err != nil {
		log.Error(err)
		return err
	}

	tripID, _ := result.LastInsertId()
	orderTrip.TripID = tripID

	// write to Cache
	data, err := json.Marshal(orderTrip)
	if err != nil {
		return err
	}
	// use orderID as the key
	repo.redis.setBytes(fmt.Sprintf(flags.CacheOrderTripPrefix, orderTrip.OrderID), []byte(data), flags.CACHETTL)
	return nil
}

func (repo *orderRepository) UpdateOrderDispacherStatus(order models.OrderDispacherStatus) error {
	t := time.Now()

	_, err := updatingDispacherStatusStmt.Exec(order.DriverID, order.Status, t, order.OrderID)
	if err != nil {
		log.Error(err)
		return err
	}

	// write to Cache
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(flags.CacheOrderPrefix, order.OrderID)
	repo.redis.delCache(cacheKey)

	orderData, _ := repo.GetOrderByID(order.OrderID)
	if orderData.OrderID != 0 {
		repo.redis.setBytes(cacheKey, []byte(data), flags.CACHETTL)
	}
	return nil
}

func (repo *orderRepository) UpdateOrder(order models.Order) error {
	t := time.Now()
	_, err := updatingOrderStmt.Exec(order.PickupTime, order.CustomerID, order.DriverID, order.PickupLongitude, order.PickupLatitude, order.DropoffLongitude, order.DropoffLatitude, order.Status, t, order.OrderID)
	if err != nil {
		log.Error(err)
		return err
	}

	//delete cache
	repo.redis.delCache(fmt.Sprintf(flags.CacheOrderPrefix, order.OrderID))
	return nil
}

func (repo *orderRepository) WriteOrder(order models.Order) (*models.Order, error) {
	t := time.Now()
	result, err := creatingOrderStmt.Exec(order.PickupTime, order.CustomerID, order.DriverID, order.PickupLongitude, order.PickupLatitude, order.DropoffLongitude, order.DropoffLatitude, flags.Prebook, t, t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	orderID, _ := result.LastInsertId()
	order.OrderID = orderID

	// write to Cache
	data, err := json.Marshal(order)
	if err == nil {
		repo.redis.setBytes(fmt.Sprintf(flags.CacheOrderPrefix, orderID), []byte(data), flags.CACHETTL)
	} else {
		log.Error(err)
	}

	return &order, nil
}

func (repo *orderRepository) GetOrderByStatus(status int) (*models.Order, error) {
	order := &models.Order{}
	err := gettingOrderByStatusStmt.QueryRow(status).Scan(&order.OrderID, &order.PickupTime, &order.CustomerName, &order.CustomerID, &order.DriverID,
		&order.PickupLongitude, &order.PickupLatitude, &order.DropoffLongitude, &order.DropoffLatitude, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return order, nil
}

func (repo *orderRepository) GetOrderByID(orderID int64) (*models.Order, error) {
	order := &models.Order{}

	cacheKey := fmt.Sprintf(flags.CacheOrderPrefix, orderID)
	cacheData, err := repo.redis.getBytes(cacheKey)
	if err == nil {
		err := json.Unmarshal(cacheData, order)
		if err == nil {
			return order, nil
		} else {
			log.Error(err)
		}
	}
	gettingOrderByIDStmt.QueryRow(orderID).Scan(&order.OrderID, &order.PickupTime, &order.CustomerID, &order.DriverID, &order.PickupLongitude, &order.PickupLatitude, &order.DropoffLongitude, &order.DropoffLatitude, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data, err := json.Marshal(order)
	if err == nil {
		repo.redis.setBytes(flags.CacheOrderPrefix+strconv.FormatInt(orderID, 10), []byte(data), flags.CACHETTL)
	} else {
		log.Error(err)
	}

	return order, nil
}
