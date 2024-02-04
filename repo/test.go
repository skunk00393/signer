package repo

import (
	"database/sql"
	"time"
)

type TestRepo struct {
	db *sql.DB
}

func NewTestRepo(db *sql.DB) TestRepo {
	return TestRepo{db: db}
}

type Signature struct {
	SignatureID string
	Signature   string
	Timestamp   time.Time
}

func (ur *TestRepo) SignTest(username string, hashStr string) (int64, error) {
	sqlStatement := `INSERT INTO signature (
			user_id,
			hashed_sig, 
			signed_on
		) 
		SELECT 
			u.id,
			$2,
			$3
		FROM 
			users u 
		WHERE 
			u.username = $1 
		RETURNING id;`
	row := ur.db.QueryRow(sqlStatement, username, hashStr, time.Now())
	var id int64
	if err := row.Err(); err != nil {
		return 0, err
	}
	row.Scan(&id)
	return id, nil
}

func (ur *TestRepo) GetSignature(username string) (*Signature, error) {
	sqlStatement := `select 
					s.id, 
					s.hashed_sig,
					s.signed_on
					from signature s 
					where 
					s.user_id = (
						select u.id 
						from users u 
						where u.username = $1)`
	row := ur.db.QueryRow(sqlStatement, username)
	sig := Signature{}
	if err := row.Err(); err != nil {
		return nil, err
	}
	row.Scan(&sig.SignatureID, &sig.Signature, &sig.Timestamp)
	return &sig, nil
}
