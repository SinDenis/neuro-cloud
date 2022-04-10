package repository

import (
	"context"
	"demo-rest/internal/domain"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertUserQuery     = "insert into \"user\" values (nextval('user_seq'), $1, $2, $3)"
	selectUserByIdQuery = "select u.id, u.username, u.password from \"user\" u where u.username = $1"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pgxPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pgxPool}
}

func (r *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	rows, err := r.pool.Query(context.Background(), selectUserByIdQuery, username)
	if err != nil {
		return domain.User{}, err
	}

	defer rows.Close()
	if !rows.Next() {
		return domain.User{}, errors.New("Not found user by username = " + username)
	}

	user := domain.User{}
	err = rows.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) SaveUser(user domain.User) error {
	rows, err := r.pool.Query(context.Background(), insertUserQuery, user.Username, user.Password, user.Email)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}
