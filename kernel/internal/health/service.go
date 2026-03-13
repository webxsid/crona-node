package health

import "context"

type Service struct {
	dbPing func(context.Context) error
}

func NewService(dbPing func(context.Context) error) *Service {
	return &Service{dbPing: dbPing}
}

type Status struct {
	DB bool `json:"db"`
}

func (s *Service) Check(ctx context.Context) Status {
	if s == nil || s.dbPing == nil {
		return Status{DB: false}
	}
	return Status{DB: s.dbPing(ctx) == nil}
}
