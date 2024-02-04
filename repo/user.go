package repo

import (
	"database/sql"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (ur *UserRepo) CreateNewUser(username string, password string) error {
	sqlStatement := `insert into users("username","password", "created_on") values($1,$2,$3)`
	_, err := ur.db.Exec(sqlStatement, username, password, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUserPassword(username string) (string, error) {
	sqlStatement := `select u.password from user u where u.username = $1`
	row := ur.db.QueryRow(sqlStatement, username)
	var pass string
	err := row.Scan(&pass)
	if err != nil {
		return "", err
	}
	return pass, nil
}
