package trip

import (
	"context"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/nkyizbay/ticket_store/internal/user"
	"gorm.io/gorm"
)

var (
	ErrDuplicateIdx = errors.New(`ERROR: duplicate key value violates unique constraint "idx_trips_idx_member" (SQLSTATE 23505)`)
	ErrTripNotFound = errors.New("this trip is not available")
)

type Repository interface {
	Create(ctx context.Context, trip *Trip) error
	Delete(ctx context.Context, id int) error
}

type defaultRepository struct {
	database *gorm.DB
}

func NewTripRepository(database *gorm.DB) Repository {
	return &defaultRepository{database: database}
}

func (t *defaultRepository) Create(ctx context.Context, trip *Trip) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := t.database.WithContext(timeoutCtx).Model(&Trip{}).Create(trip).Error; err != nil {
		if err.Error() == ErrDuplicateIdx.Error() {
			return ErrDuplicateIdx
		}

		log.Error(err)
		return err
	}

	return nil
}

func (t *defaultRepository) Delete(ctx context.Context, id int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := t.database.WithContext(timeoutCtx).Delete(&Trip{}, id).Error; err != nil {
		switch {
		case errors.Is(err, user.ErrNoRecord):
			return ErrTripNotFound
		default:
			log.Error(err)
			return err
		}
	}

	return nil
}
