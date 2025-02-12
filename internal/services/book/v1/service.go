package bookv1

import "github.com/FotiadisM/service-template/internal/services/book/v1/queries"

type Service struct {
	db queries.Querier
}

func NewService(db queries.Querier) *Service {
	return &Service{db: db}
}
