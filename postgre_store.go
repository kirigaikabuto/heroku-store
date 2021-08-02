package heroku_store

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	_ "github.com/lib/pq"
)

var userPostgresQueries = []string{
	`create table if not exists users(
		id text,
		username text,
		primary key(id)
	);`,
}

type userStore struct {
	db *sql.DB
}

func NewUsersPostgreStore(cfg PostgresConfig) (UserStore, error) {
	db, err := getDbConn(getConnString(cfg))
	if err != nil {
		return nil, err
	}
	for _, q := range userPostgresQueries {
		_, err := db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}

	db.SetMaxOpenConns(10)
	store := &userStore{db: db}
	store.Create(&User{"1", "kirito"})
	return store, nil
}

func (u *userStore) Create(user *User) (*User, error) {
	user.Id = uuid.New().String()
	query := "insert into users (id, username) values ($1, $2)"
	result, err := u.db.Exec(query, user.Id, user.Username)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, errors.New("user creating error")
	}
	return user, nil
}

func (u *userStore) List() ([]User, error) {
	items := []User{}
	query := "select id, username from users "

	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := User{}
		err = rows.Scan(&item.Id, &item.Username)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
