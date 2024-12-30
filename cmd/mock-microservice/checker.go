package main

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"

	"github.com/FotiadisM/mock-microservice/internal/db"
)

var errUnknownService = errors.New("unknown service")

type checker struct {
	DB db.DB
}

var _ grpchealth.Checker = &checker{}

func (c *checker) Check(ctx context.Context, req *grpchealth.CheckRequest) (*grpchealth.CheckResponse, error) {
	if req.Service == "startup" {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
	}

	if req.Service == "readiness" {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
	}

	if req.Service == "liveness" {
		return c.livenessProbe(ctx)
	}

	err := connect.NewError(
		connect.CodeNotFound,
		fmt.Errorf("'%s': %w", req.Service, errUnknownService),
	)

	return nil, err
}

func (c *checker) livenessProbe(ctx context.Context) (*grpchealth.CheckResponse, error) {
	err := c.DB.Ping(ctx)
	if err != nil {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusNotServing}, nil //nolint
	}

	return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
}
