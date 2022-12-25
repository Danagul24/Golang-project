package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"project/internal/models"
	"project/internal/store"
)

func (db *DB) Users() store.UsersRepository {
	if db.users == nil {
		db.users = newUserRepository(db.conn)
	}
	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func newUserRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (u UsersRepository) Create(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if err := user.BeforeCreating(); err != nil {
		return err
	}
	_, err := u.conn.Exec("INSERT INTO users(name, surname, email, password, phone_number, birth_date, role) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Name, user.Surname, user.Email, user.EncryptedPassword, user.PhoneNumber, user.BirthDate, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u UsersRepository) All(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)
	basicQuery := "SELECT * FROM users"

	if err := u.conn.Select(&users, basicQuery); err != nil {
		return nil, err
	}
	return users, nil
}

func (u UsersRepository) ByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.Get(user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}
	return user, nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if err := user.BeforeCreating(); err != nil {
		return err
	}
	_, err := u.conn.Exec("UPDATE users SET name = $1, surname = $2, password = $3, phone_number = $4, birth_date = $5, role = $6",
		user.Name, user.Surname, user.Password, user.PhoneNumber, user.BirthDate, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id int) error {
	if _, err := u.conn.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}
