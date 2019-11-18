package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"time"

	"github.com/interview-order/config"
	"github.com/interview-order/endpoint"
	"github.com/interview-order/pb"
	"github.com/interview-order/repository"
	"github.com/interview-order/server"
	"github.com/interview-order/service"
	"github.com/interview-order/util"
	log "github.com/sirupsen/logrus"
)

func main() {

	//init config
	configFile := flag.String("c", "", "Configuration File")
	flag.Parse()

	if *configFile == "" {
		log.Info("\n\nUse -h to get more information on command line options\n")
		log.Info("You must specify a configuration file")
		os.Exit(1)
	}

	err := config.Initialize(*configFile)
	if err != nil {
		log.Printf("Error reading configuration: %s\n", err.Error())
		os.Exit(1)
	}

	err = repository.InitDbFactory()
	if err != nil {
		log.Println("Cannot connect to database: ", err.Error())
		return
	}

	util.InitTimeZoneLocation()

	start := time.Now()

	httpAddr := fmt.Sprintf(":%d", config.MustGetInt("server.http_port"))
	gRPCAddr := fmt.Sprintf(":%d", config.MustGetInt("server.grpc_port"))

	setupLogging()
	log.Info("Server Started in ", time.Since(start))
	defer log.Info("Closed")

	ctx := context.Background()

	// init lorem service
	var svcPing service.IPingService
	svcPing = service.PingService{}

	var svcOrder service.OrderServiceInterface
	svcOrder = service.OrderService{}

	// creating Endpoints struct
	endpoints := endpoint.Endpoints{
		PingEndpoint:                       endpoint.MakePingEndpoint(svcPing),
		GetOrderByIDEndpoint:               endpoint.MakeGetOrderByIDEndpoint(svcOrder),
		CreateOrderEndpoint:                endpoint.MakeCreateOrderEndpoint(svcOrder),
		GetTripByOrderIDEndpoint:           endpoint.MakeGetTripByOrderIDEndpoint(svcOrder),
		UpdateOrderDispacherStatusEndpoint: endpoint.MakeUpdateOrderDispacherStatusEndpoint(svcOrder),
	}

	// initiate pubsub service
	service.InitMQPublisher()
	service.InitMQConsumer()

	errChan := make(chan error)

	// Error channel.
	errc := make(chan error)

	// HTTP transport.
	go func() {
		log.Info("transport HTTP addr", httpAddr)

		handler := server.NewHttpServer(ctx, endpoints)
		errc <- http.ListenAndServe(httpAddr, handler)
	}()

	// gRPC transport.
	go func() {
		log.Info("transport gRPC addr", gRPCAddr)

		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := server.NewGRPCServer(ctx, endpoints)
		s := grpc.NewServer()
		pb.RegisterPingServer(s, srv)
		errc <- s.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	log.Info(<-errChan)
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	if config.MustGetString("server.mode") == "production" {
		log.Info("here")
		logPath := config.MustGetString("server.log_path")

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			log.Fatal("Cannot log to file", err.Error())
		}

		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(file)
	}
}
