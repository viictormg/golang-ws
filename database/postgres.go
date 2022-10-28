package database

import (
	"context"
	"database/sql"
	"fmt"
	"ws-go/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	sql := "INSERT INTO users (id, email, password) VALUES($1, $2, $3)"
	_, err := repo.db.ExecContext(ctx, sql, user.Id, user.Email, user.Password)

	return err

}

func (repo *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	fmt.Println(id)
	sql := "SELECT id, email FROM users WHERE id = $1"

	rows, err := repo.db.QueryContext(ctx, sql, id)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		// log.Fatal(err)
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	sql := "SELECT id, email, password FROM users WHERE email = $1"

	rows, err := repo.db.QueryContext(ctx, sql, email)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		// log.Fatal(err)
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
