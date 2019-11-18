# Order Service For Interview


### Pre-requirement to run the services
*  1. Redis          (v. 4.0.9+)
*  2. Mysql       (v. 5.7.27)

##  Generate pb file
*   grpc
*   protoc --go_out=plugins=grpc:. pb/*.proto

## Config
* config.toml

## Start local
* go run main.go -config config.toml
