# Order Service For Interview


### Pre-requirement to run the services
*  1. Redis          (v. 4.0.9+)
*  2. PostgreSQL     (v. 11.3)
*  3. RabbitMQ       (v. 3.7.8+)    

##  Generate pb file
*   grpc
*   protoc --go_out=plugins=grpc:. pb/*.proto

## Config
* config.toml

## Start local
* go run main.go -c config.toml
