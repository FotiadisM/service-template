// Package health provies an in-process grpc_health_v1.HealthClient
// to a grpc_health_v1.HealthServer
package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Client is an in-process grpc_health_v1.Client to a grpc_health_v1.HealthServer.
type Client struct {
	svc grpc_health_v1.HealthServer

	grpc_health_v1.HealthClient
}

func NewHealthClient(svc grpc_health_v1.HealthServer) *Client {
	return &Client{
		svc: svc,
	}
}

func (hc *Client) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, _ ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error) {
	return hc.svc.Check(ctx, in)
}
