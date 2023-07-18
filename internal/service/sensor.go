package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/internal/domain"
	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/cache"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
)

var ErrorWrongGroupName error = errors.New("wrong group name")

type SensorService interface {
	GetTransparency(ctx context.Context, groupName string) (*float64, error)
	GetTemperature(ctx context.Context, groupName string) (*float64, error)
	GetCurrentSpecies(ctx context.Context, groupName string) ([]domain.DetectedFish, error)
	GetCurrentTopSpecies(ctx context.Context, groupName, start, end string, top int) ([]domain.DetectedFish, error)
	GetRegionTemperature(ctx context.Context, region domain.Region, flag string) (float64, error)
	GetSensorTemperature(ctx context.Context, inGroupID int, group, start, end string) (*float64, error)
}

type Service struct {
	db    storage.SensorPostgres
	log   logging.Logger
	ctx   context.Context
	cfg   config.Config
	cache cache.CacheRedis
}

func NewService(ctx context.Context, db storage.SensorPostgres, log logging.Logger, cfg config.Config, cache cache.CacheRedis) *Service {
	return &Service{db: db, log: log, ctx: ctx, cfg: cfg, cache: cache}
}

func (s *Service) GetTransparency(ctx context.Context, groupName string) (*float64, error) {
	if !s.validateGroupName(groupName) {
		return nil, ErrorWrongGroupName
	}

	cacheKey := fmt.Sprintf("transparency for %s", groupName)

	exist, err := s.cache.IfExistsInCache(ctx, cacheKey)
	if err != nil {
		s.log.Error(cache.ErrorCheckExist)
		return nil, err
	}

	if !exist {
		transparency, err := s.db.GetTransparency(ctx, groupName)
		if err != nil {
			return nil, fmt.Errorf("error get transparency from DB, err: %v", err)
		}

		err = s.cache.Set(ctx, cacheKey, strconv.FormatFloat(transparency, 'f', -1, 32))
		if err != nil {
			s.log.Error(cache.ErrorSetRedis)
			return nil, err
		}
		return &transparency, nil
	} else {
		cachedTr, err := s.cache.Get(ctx, cacheKey)
		if err != nil {
			s.log.Error(cache.ErrorGetRedis)
			return nil, err
		}

		floatTr, _ := strconv.ParseFloat(cachedTr, 64)

		s.log.Info("returned transparency from cache")
		return &floatTr, nil
	}
}

func (s *Service) GetTemperature(ctx context.Context, groupName string) (*float64, error) {
	if !s.validateGroupName(groupName) {
		return nil, ErrorWrongGroupName
	}

	cacheKey := fmt.Sprintf("temperature for %s", groupName)

	exist, err := s.cache.IfExistsInCache(ctx, cacheKey)
	if err != nil {
		s.log.Error(cache.ErrorCheckExist)
		return nil, err
	}

	if !exist {
		temperature, err := s.db.GetTemperature(ctx, groupName)
		if err != nil {
			return nil, fmt.Errorf("error get temperature from DB, err: %v", err)
		}

		err = s.cache.Set(ctx, cacheKey, strconv.FormatFloat(temperature, 'f', -1, 32))
		if err != nil {
			s.log.Error(cache.ErrorSetRedis)
			return nil, err
		}
		return &temperature, nil
	} else {
		cachedTmp, err := s.cache.Get(ctx, cacheKey)
		if err != nil {
			s.log.Error(cache.ErrorGetRedis)
			return nil, err
		}

		floatTmp, _ := strconv.ParseFloat(cachedTmp, 64)

		s.log.Info("returned temperature from cache")
		return &floatTmp, nil
	}
}

func (s *Service) GetCurrentSpecies(ctx context.Context, groupName string) ([]domain.DetectedFish, error) {
	if !s.validateGroupName(groupName) {
		return nil, ErrorWrongGroupName
	}

	species, err := s.db.GetSpecies(ctx, groupName)
	if err != nil {
		return nil, err
	}

	return species, nil
}

func (s *Service) GetCurrentTopSpecies(ctx context.Context, groupName, start, end string, top int) ([]domain.DetectedFish, error) {
	if !s.validateGroupName(groupName) {
		return nil, ErrorWrongGroupName
	}

	species, err := s.db.GetTopSpecies(ctx, groupName, start, end, top)
	if err != nil {
		return nil, err
	}

	return species, nil
}

func (s *Service) GetRegionTemperature(ctx context.Context, region domain.Region, flag string) (float64, error) {
	temperature, err := s.db.GetRegionTemperature(ctx, region, flag)
	if err != nil {
		return 0, err
	}

	return temperature, nil
}

func (s *Service) GetSensorTemperature(ctx context.Context, inGroupID int, group, start, end string) (*float64, error) {
	temperature, err := s.db.GetSensorAverageTemperature(ctx, inGroupID, group, start, end)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}

func (s *Service) validateGroupName(groupName string) bool {
	exist := false
	for _, name := range strings.Split(s.cfg.GroupNames, " ") {
		if groupName == name {
			exist = true
		}
	}

	if exist {
		return true
	}

	return false
}
