package service

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
)

// Define service interface
type IPingService interface {
	// generate a word with at least min letters and at most max letters.
	Ping(ctx context.Context, req *empty.Empty) (string, error)
}

// Implement service with empty struct
type PingService struct {
}

// Implement service functions
func (PingService) Ping(ctx context.Context, req *empty.Empty) (string, error) {
	return "pong", nil
}
