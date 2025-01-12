package authv1

import (
	"github.com/FotiadisM/mock-microservice/internal/services/auth/v1/queries"
)

type Service struct {
	db queries.Querier
}

func NewService(db queries.Querier) *Service {
	return &Service{db: db}
}
