package models

type GeneraResponse struct {
	Response string `json:"Response"`
}

type Order struct {
	OrderID          int64   `json:"order_id"`
	CustomerID       int64   `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	CustomerMobile   string  `json:"customer_mobile"`
	DriverID         int64   `json:"driver_id"`
	DriverName       string  `json:"driver_name"`
	DriverMobile     string  `json:"driver_mobile"`
	VehicleID        string  `json:"vehicle_id"`
	VehicleNo        string  `json:"vehicle_no"`
	PickupTime       string  `json:"pick_time"`
	PickupLatitude   float64 `json:"pickup_latitude"`
	PickupLongitude  float64 `json:"pickup_longitude"`
	DropoffLatitude  float64 `json:"dropoff_latitude"`
	DropoffLongitude float64 `json:"dropoff_longitude"`
	Status           int64   `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type OrderTrip struct {
	TripID           int64   `json:"trip_id"`
	OrderID          int64   `json:"order_id"`
	CustomerID       int64   `json:"customer_id"`
	DriverID         int64   `json:"driver_id"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	VacantDistance   float64 `json:"vacant_distance"`
	EngagedDistance  float64 `json:"engaged_distance"`
	PickupLatitude   float64 `json:"pickup_latitude"`
	PickupLongitude  float64 `json:"pickup_longitude"`
	DropoffLatitude  float64 `json:"dropoff_latitude"`
	DropoffLongitude float64 `json:"dropoff_longitude"`
	Status           int64   `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type OrderDispacherStatus struct {
	OrderID         int64   `json:"order_id"`
	DriverID        int64   `json:"driver_id"`
	Status          int64   `json:"status"`
	DriverLatitude  float64 `json:"driver_latitude"`
	DriverLongitude float64 `json:"driver_longitude"`
}

type Driver struct {
	DriverID     int64   `json:"driver_id"`
	OrderID      int64   `json:"order_id"`
	DriverName   string  `json:"driver_name"`
	DriverMobile string  `json:"driver_mobile"`
	VehicleID    string  `json:"vehicle_id"`
	VehicleNo    string  `json:"vehicle_no"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Status       int     `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type Passenger struct {
	PassengerID     string `json:"passengerId"`
	PassengerName   string `json:"passengerName"`
	PassengerMobile string `json:"passengerMobile"`
	CreatedAt       string `json:"CreatedAt"`
	UpdatedAt       string `json:"UpdatedAt"`
}

type GeneratedPassengerResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Response struct {
		Type string    `json:"@type"`
		Data Passenger `json:"data"`
	} `json:"response"`
}

type IDRequest struct {
	ID int64 `json:"id"`
}

type GeneratedDriverResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Response struct {
		Type string     `json:"@type"`
		Data DriverInfo `json:"data"`
	} `json:"response"`
}

type DriverInfo struct {
	DriverID     string `json:"driverId"`
	DriverName   string `json:"driverName"`
	DriverMobile string `json:"driverMobile"`
	VehicleID    string `json:"vehicleId"`
	VehicleNo    string `json:"vehicleNo"`
	CreatedAt    string `json:"CreatedAt"`
	UpdatedAt    string `json:"UpdatedAt"`
}

type UserNotification struct {
	UserID  int64  `json:"userID"`
	Message string `json: "message"`
}
