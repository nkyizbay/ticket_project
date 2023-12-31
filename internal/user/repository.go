package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrNoRecord         = errors.New("there is no record in DB with that username")
	ErrUniqueViolation  = errors.New("UNIQUE constraint failed")
	DuplicateEntryError = &pgconn.PgError{Code: "23505"}
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByUserName(ctx context.Context, username string) (*User, error)
}

type defaultRepository struct {
	database *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &defaultRepository{
		database: database,
	}
}

func (r *defaultRepository) GetByUserName(ctx context.Context, username string) (*User, error) {
	user := User{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&User{}).
		First(&user, "user_name = ?", username).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
	}

	return &user, nil
}

func (r *defaultRepository) Create(ctx context.Context, user *User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&User{}).Create(user).Error; err != nil {
		switch {
		case errors.As(err, &DuplicateEntryError):
			return ErrDuplicatedValue
		default:
			return err
		}
	}

	return nil
}
