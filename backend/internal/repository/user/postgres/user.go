package user_postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.UserInfo) error {
	query := `
		INSERT INTO users (id, name, registration_date, birth_date, access_level)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.RegistrationDate,
		user.BirthDate,
		int(user.AccessLvl),
	)

	return err
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error) {
	query := `
		SELECT id, name, registration_date, birth_date, access_level
		FROM users
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, userID)

	var user entity.UserInfo
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.RegistrationDate,
		&user.BirthDate,
		&user.AccessLvl,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.UserInfo) error {
	query := `
		UPDATE users
		SET name = $1, birth_date = $2
		WHERE id = $3
	`
	cmdTag, err := r.pool.Exec(ctx, query,
		user.Name,
		user.BirthDate,
		user.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
