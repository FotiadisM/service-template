package filters

import (
	"strings"

	"github.com/FotiadisM/mock-microservice/pkg/otelgrpc"
)

func FullMethodName(s string) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		return s == fullMethodName
	}
}

func MethodName(s string) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		_, m := ParseFullMethodName(fullMethodName)
		return s == m
	}
}

func MethodPrefix(pre string) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		_, m := ParseFullMethodName(fullMethodName)
		return strings.HasPrefix(m, pre)
	}
}

func ServiceName(s string) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		svc, _ := ParseFullMethodName(fullMethodName)
		return s == svc
	}
}

func ServicePrfix(pre string) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		svc, _ := ParseFullMethodName(fullMethodName)
		return strings.HasPrefix(svc, pre)
	}
}

func HealthCheck() otelgrpc.Filter {
	return func(fullMethodName string) bool {
		svc, _ := ParseFullMethodName(fullMethodName)
		return svc == "grpc.health.v1.Health"
	}
}

func All(fs ...otelgrpc.Filter) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		for _, f := range fs {
			if !f(fullMethodName) {
				return false
			}
		}
		return true
	}
}

func Any(fs ...otelgrpc.Filter) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		for _, f := range fs {
			if f(fullMethodName) {
				return true
			}
		}
		return false
	}
}

func None(fs ...otelgrpc.Filter) otelgrpc.Filter {
	return Not(Any(fs...))
}

func Not(f otelgrpc.Filter) otelgrpc.Filter {
	return func(fullMethodName string) bool {
		return !f(fullMethodName)
	}
}

func ParseFullMethodName(s string) (service, method string) {
	name := strings.TrimLeft(s, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
