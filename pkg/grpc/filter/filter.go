package filter

import (
	"strings"
)

type Filter func(fullMethodName string) bool

func FullMethodName(s string) Filter {
	return func(fullMethodName string) bool {
		return s == fullMethodName
	}
}

func MethodName(s string) Filter {
	return func(fullMethodName string) bool {
		_, m := parseFullMethodName(fullMethodName)
		return s == m
	}
}

func MethodPrefix(pre string) Filter {
	return func(fullMethodName string) bool {
		_, m := parseFullMethodName(fullMethodName)
		return strings.HasPrefix(m, pre)
	}
}

func ServiceName(s string) Filter {
	return func(fullMethodName string) bool {
		svc, _ := parseFullMethodName(fullMethodName)
		return s == svc
	}
}

func ServicePrfix(pre string) Filter {
	return func(fullMethodName string) bool {
		svc, _ := parseFullMethodName(fullMethodName)
		return strings.HasPrefix(svc, pre)
	}
}

func HealthCheck() Filter {
	return func(fullMethodName string) bool {
		svc, _ := parseFullMethodName(fullMethodName)
		return svc == "grpc.health.v1.Health"
	}
}

func Reflection() Filter {
	return func(fullMethodName string) bool {
		svc, _ := parseFullMethodName(fullMethodName)
		return (svc == "grpc.reflection.v1.ServerReflection" || svc == "grpc.reflection.v1alpha.ServerReflection")
	}
}

func All(fs ...Filter) Filter {
	return func(fullMethodName string) bool {
		for _, f := range fs {
			if !f(fullMethodName) {
				return false
			}
		}
		return true
	}
}

func Any(fs ...Filter) Filter {
	return func(fullMethodName string) bool {
		for _, f := range fs {
			if f(fullMethodName) {
				return true
			}
		}
		return false
	}
}

func None(fs ...Filter) Filter {
	return Not(Any(fs...))
}

func Not(f Filter) Filter {
	return func(fullMethodName string) bool {
		return !f(fullMethodName)
	}
}

func parseFullMethodName(s string) (service, method string) {
	parts := strings.SplitN(strings.TrimLeft(s, "/"), "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
