package database

import (
	"database/sql"
	"errors"

	"github.com/avalchev94/sqlxt"
	"github.com/lib/pq"
)

// Task represent a single task that a user can add
type Task struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id" sql:"user_id"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Labels Labels `json:"labels,omitempty" sql:"labels"`
}

var (
	InvalidUser = errors.New("db: invalid user_id")
	EmptyBody   = errors.New("db: body can't be empty")
)

// OK validates the task fields.
func (t *Task) OK() error {
	switch {
	case t.UserID <= 0:
		return InvalidUser
	case len(t.Body) == 0:
		return EmptyBody
	}
	return nil
}

func (t *Task) add(tx *sql.Tx) error {
	if err := t.OK(); err != nil {
		return err
	}

	row := tx.QueryRow("INSERT INTO tasks (user_id, title, body, labels) VALUES($1,$2,$3,$4) RETURNING id",
		t.UserID, t.Title, t.Body, pq.Array(t.Labels))

	return row.Scan(&t.ID)
}

// GetTask searches for task with the specified id
func GetTask(taskID int64, db *sql.DB) (*Task, error) {
	var t Task

	scanner := sqlxt.NewScanner(db.Query("SELECT * FROM tasks WHERE id=$1", taskID))
	if err := scanner.Scan(&t); err != nil {
		return nil, err
	}

	return &t, nil
}
