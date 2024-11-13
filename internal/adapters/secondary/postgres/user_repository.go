package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/ncfex/dcart-auth/internal/adapters/secondary/postgres/db"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
	"github.com/ncfex/dcart-auth/internal/core/ports/outbound"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(database *database) outbound.UserRepository {
	return &userRepository{
		queries: db.New(database.DB),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, userObj *userDomain.User) (*userDomain.User, error) {
	params := db.CreateUserParams{
		Username:     userObj.Username,
		PasswordHash: userObj.PasswordHash,
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserAlreadyExists
		}
		return nil, err
	}

	return db.ToUserDomain(&dbUser), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID *uuid.UUID) (*userDomain.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, *userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return db.ToUserDomain(&dbUser), nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*userDomain.User, error) {
	dbUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return db.ToUserDomain(&dbUser), nil
}
