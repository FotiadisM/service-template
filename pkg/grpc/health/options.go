package health

import (
	"context"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type ProbeFunc func(context.Context) (healthv1.HealthCheckResponse_ServingStatus, error)

func defaultProbeFunc(_ context.Context) (healthv1.HealthCheckResponse_ServingStatus, error) {
	return healthv1.HealthCheckResponse_SERVING, nil
}

type options struct {
	readiness ProbeFunc
	liveness  ProbeFunc
	startup   ProbeFunc
}

func defaultOptions() *options {
	return &options{
		readiness: defaultProbeFunc,
		liveness:  defaultProbeFunc,
		startup:   defaultProbeFunc,
	}
}

type Option func(o *options)

func WithReadiness(fn ProbeFunc) Option {
	return func(o *options) {
		o.readiness = fn
	}
}

func WithLiveness(fn ProbeFunc) Option {
	return func(o *options) {
		o.readiness = fn
	}
}

func WithStartup(fn ProbeFunc) Option {
	return func(o *options) {
		o.readiness = fn
	}
}
