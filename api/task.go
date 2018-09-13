package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Task represent a single task that a user can add
type Task struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Title    string    `json:"title,omitempty"`
	Data     string    `json:"data"`
	Labels   []Label   `json:"labels,omitempty"`
	LabelsID []int64   `json:"labels_id,omitempty"`
	Created  time.Time `json:"created"`
}

// OK validates the task fields.
func (t *Task) OK() error {
	if t.UserID <= 0 {
		return errors.New("UserID is incorrect")
	}
	if len(t.Data) == 0 {
		return errors.New("task description can't be empty")
	}
	return nil
}

// Add creates a task in the database.
func (t *Task) Add(db *sql.DB) error {
	if err := t.OK(); err != nil {
		return err
	}

	t.Created = time.Now()
	_, err := db.Exec("INSERT INTO tasks (user_id, title, data, labels, created) VALUES($1,$2,$3,$4,$5)",
		t.UserID, t.Title, t.Data, pq.Array(t.LabelsID), t.Created)
	return err
}

func (t *Task) getLabels(db *sql.DB) error {
	if err := t.OK(); err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, name, color FROM labels WHERE id=ANY($1)", pq.Array(t.LabelsID))
	if err != nil {
		return err
	}

	for rows.Next() {
		var label Label
		if err := rows.Scan(&label.ID, &label.Name, &label.Color); err != nil {
			break
		}
		t.Labels = append(t.Labels, label)
	}
	return nil
}

// GetTask searches for task with the specified id
func GetTask(taskID int64, db *sql.DB) (*Task, error) {
	row := db.QueryRow("SELECT * FROM tasks WHERE id=$1", taskID)

	var t Task
	err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Data, pq.Array(&t.LabelsID), &t.Created)
	if err != nil {
		return nil, err
	}
	t.getLabels(DB)

	return &t, nil
}

// GetTasks returns all the tasks for some user.
func GetTasks(userID int64, db *sql.DB) ([]*Task, error) {
	rows, err := db.Query("SELECT * FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Data, pq.Array(&t.LabelsID), &t.Created)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	// get the labels for each task
	for _, task := range tasks {
		task.getLabels(DB)
	}

	return tasks, nil
}
