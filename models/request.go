package models

import (
	"database/sql"
	"log"
)

//Tasks table
type Tasks struct {
	ID        int `json:"id"`
	checkerID int
	userID    int
	Interval  int    `json:"interval"`
	Target    string `json:"target"`
	slackID   int
	status    bool
}

//Checker table
type Checker struct {
	CheckerID int `json:"id"`
	hashID    int
	status    bool
}

//FailOnRequest func
func FailOnRequest(err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//CreateTables one transaction
func CreateTables(db *sql.DB) error {
	tx, err := db.Begin()
	FailOnRequest(err)
	defer tx.Rollback()

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS hash (hash_id SERIAL PRIMARY KEY, hash VARCHAR UNIQUE NOT NULL)")
	FailOnRequest(err)
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS users (user_id SERIAL PRIMARY KEY, login VARCHAR UNIQUE NOT NULL, password VARCHAR NOT NULL, hash_id INTEGER UNIQUE NOT NULL REFERENCES hash(hash_id), hash_activate BOOL NOT NULL DEFAULT false)")
	FailOnRequest(err)
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS slack (slack_id SERIAL PRIMARY KEY, user_id INTEGER NOT NULL REFERENCES users(user_id), url VARCHAR NOT NULL, channel VARCHAR NOT NULL)")
	FailOnRequest(err)
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS checkers (checker_id SERIAL PRIMARY KEY, hash_id INTEGER UNIQUE NOT NULL REFERENCES hash(hash_id), status BOOL NOT NULL)")
	FailOnRequest(err)
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS servers (server_id SERIAL PRIMARY KEY, checker_id INTEGER NOT NULL REFERENCES checkers(checker_id), timepointer timestamp NOT NULL)")
	FailOnRequest(err)
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS tasks (task_id SERIAL PRIMARY KEY, checker_id INTEGER NOT NULL REFERENCES checkers(checker_id), user_id INTEGER NOT NULL REFERENCES users(user_id), interval INTEGER NOT NULL, target VARCHAR NOT NULL, slack_id INTEGER NOT NULL REFERENCES slack(slack_id), status BOOL NOT NULL DEFAULT true)")
	FailOnRequest(err)

	err = tx.Commit()
	FailOnRequest(err)

	return nil
}

//GetTasksReq SELECT all task from checker
func GetTasksReq(db *sql.DB, id int) ([]*Tasks, error) {
	rows, err := db.Query("SELECT * FROM tasks WHERE checker_id = $1", id)
	FailOnRequest(err)
	defer rows.Close()

	bks := make([]*Tasks, 0)
	for rows.Next() {
		bk := new(Tasks)
		err = rows.Scan(&bk.ID, &bk.checkerID, &bk.userID, &bk.Interval, &bk.Target, &bk.slackID, &bk.status)
		FailOnRequest(err)
		bks = append(bks, bk)
	}
	err = rows.Err()
	FailOnRequest(err)

	return bks, nil
}

//GetCheckerID SELECT checker_id with hash_id
func GetCheckerID(db *sql.DB, hash string) (Checker, error) {
	stmt, err := db.Prepare("SELECT * FROM checkers WHERE hash_id=(SELECT hash_id from hash WHERE hash=$1)")
	FailOnRequest(err)

	var id Checker
	err = stmt.QueryRow(hash).Scan(&id.CheckerID, &id.hashID, &id.status)
	FailOnRequest(err)

	return id, nil
}

//InsertHash INSERT new clients
func InsertHash(db *sql.DB, hash string) error {
	query := `INSERT INTO checkers (hash_id, status)
            VALUES ((SELECT hash_id from hash WHERE hash=$1), true)
            ON CONFLICT DO NOTHING;`

	stmt, err := db.Prepare(query)
	FailOnRequest(err)

	_, err = stmt.Exec(hash)
	FailOnRequest(err)

	return nil
}

//UpdateStatus Update targer status
func UpdateStatus(db *sql.DB, id int, status bool) error {
	stmt, err := db.Prepare("UPDATE tasks SET status = $2 WHERE task_id = $1;")
	FailOnRequest(err)

	_, err = stmt.Exec(id, status)
	FailOnRequest(err)

	return nil
}
