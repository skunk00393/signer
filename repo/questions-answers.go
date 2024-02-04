package repo

import (
	"database/sql"
	"fmt"
	"time"
)

type QuastionAnswer struct {
	Question string
	Answer   string
}

type QA struct {
	db *sql.DB
}

func NewQARepo(db *sql.DB) QA {
	return QA{db: db}
}

func (qa *QA) InsertQAs(question string, answer string, id int64) error {
	sqlStatement := `insert into qas(question,answer, answered_on, sig_id) values($1,$2,$3,$4)`
	_, err := qa.db.Exec(sqlStatement, question, answer, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (qa *QA) GetQAs(signature string) (string, error) {
	sqlStatement := `select qa.question, qa.answer from qas qa where qa.sig_id = $1`
	rows, err := qa.db.Query(sqlStatement, signature)
	if err != nil {
		return "", err
	}
	var QAs string

	for rows.Next() {
		QA := QuastionAnswer{}
		err := rows.Scan(&QA.Question, &QA.Answer)
		fmt.Print(QA)

		if err != nil {
			return "", err
		}
		QAs += fmt.Sprintf("%s: %s \n", QA.Question, QA.Answer)
	}
	return QAs, nil
}
