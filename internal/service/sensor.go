package service

import (
	"context"
	"fmt"

	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
)

type SensorService interface {
	GetTransparency(ctx context.Context) error
}

type Service struct {
	db  storage.SensorPostgres
	log logging.Logger
	ctx context.Context
}

func (s *Service) GetTransparency(ctx context.Context) error {
	transparency, err := s.db.GetTransparency(ctx)
	if err != nil {
		return fmt.Errorf("error get transparency from DB, err: %v", err)
	}
}
