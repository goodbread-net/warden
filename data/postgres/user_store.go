package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/matthiase/warden/models"
	"github.com/segmentio/ksuid"
)

type UserStore struct {
	sqlx.Ext
}

func (db *UserStore) Create(firstName string, lastName string, email string) (*models.User, error) {
	sql := `
		INSERT INTO users (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, first_name, last_name, email
	`

	uid := ksuid.New()
	result, err := db.Queryx(sql, uid.String(), firstName, lastName, email)

	if err != nil {

		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, errors.New(models.ErrUserDuplicateEmail)
		}
	}

	//if err, ok := err.(*pg.Error); ok && err.Code == "23505" {
	//	return nil, models.ErrDuplicateEmail
	//	return nil, err
	//}

	defer result.Close()
	result.Next()

	user := models.User{}
	if err := result.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *UserStore) Find(id string) (*models.User, error) {
	sql := `
		SELECT id, first_name, last_name, email
		FROM users
		WHERE id = $1
	`

	user := models.User{}
	if err := sqlx.Get(db, &user, sql, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *UserStore) FindByEmail(email string) (*models.User, error) {
	sql := `
		SELECT id, first_name, last_name, email
		FROM users
		WHERE email = $1
	`

	user := models.User{}
	if err := sqlx.Get(db, &user, sql, email); err != nil {
		return nil, err
	}
	return &user, nil
}
