package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
)

var errUnknownService = errors.New("unknown service")

type healthChecker struct {
	DB *sql.DB
}

var _ grpchealth.Checker = &healthChecker{}

func (c *healthChecker) Check(ctx context.Context, req *grpchealth.CheckRequest) (*grpchealth.CheckResponse, error) {
	if req.Service == "startup" {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
	}

	if req.Service == "readiness" {
		return c.readiness(ctx)
	}

	if req.Service == "liveness" {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
	}

	err := connect.NewError(
		connect.CodeNotFound,
		fmt.Errorf("'%s': %w", req.Service, errUnknownService),
	)

	return nil, err
}

func (c *healthChecker) readiness(ctx context.Context) (*grpchealth.CheckResponse, error) {
	err := c.DB.PingContext(ctx)
	if err != nil {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusNotServing}, nil //nolint
	}

	return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
}
