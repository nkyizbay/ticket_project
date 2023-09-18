package trip

import (
	"context"
	"errors"
)

var (
	ErrAlreadyCreatedTrip = errors.New("this trip is already created")
	ErrTripNotExist       = errors.New("this trip does not exist")
)

type Service interface {
	CreateTrip(ctx context.Context, trip *Trip) error
	CancelTrip(ctx context.Context, id int) error
}

type defaultService struct {
	tripRepo Repository
}

func NewTripService(tripRepo Repository) Service {
	return &defaultService{tripRepo: tripRepo}
}

func (s *defaultService) CreateTrip(ctx context.Context, t *Trip) error {
	if err := s.tripRepo.Create(ctx, t); err != nil {
		if errors.Is(err, ErrDuplicateIdx) {
			return ErrAlreadyCreatedTrip
		}
		return err
	}

	return nil
}

func (s *defaultService) CancelTrip(ctx context.Context, id int) error {
	if err := s.tripRepo.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, ErrTripNotFound):
			return ErrTripNotExist
		}
		return err
	}

	return nil
}
