package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/avalchev94/sqlxt"

	"github.com/rickb777/date"
	"github.com/rickb777/date/clock"
)

type ScheduleType string

const (
	Unscheduled ScheduleType = "unscheduled"
	OnDate      ScheduleType = "on_date"
	DueDate     ScheduleType = "due_date"
)

func (st ScheduleType) OK() bool {
	return string(st) == string(Unscheduled) ||
		string(st) == string(OnDate) ||
		string(st) == string(DueDate)
}

type TaskSchedule struct {
	ID       int64        `json:"id,omitempty"`
	TaskID   int64        `json:"task_id,omitempty" sql:"task_id"`
	Type     ScheduleType `json:"type"`
	Date     *date.Date   `json:"date,omitempty"`
	Time     *clock.Clock `json:"time,omitempty"`
	Created  time.Time    `json:"created"`
	Finished *time.Time   `json:"finished,omitempty"`
}

func (ts *TaskSchedule) OK() error {
	switch {
	case ts.TaskID <= 0:
		return errors.New("incorrect task id")
	case !ts.Type.OK():
		return errors.New("invalid type")
	}
	return nil
}

func (ts *TaskSchedule) add(tx *sql.Tx) error {
	if err := ts.OK(); err != nil {
		return err
	}
	ts.Created = time.Now()
	//TODO: check if the task is already scheduled?
	row := tx.QueryRow(`INSERT INTO task_schedule (task_id,type,date,time,created,finished)
											VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		ts.TaskID, ts.Type, ts.Date, ts.Time, ts.Created, ts.Finished)

	return row.Scan(&ts.ID)
}

func GetSchedule(taskID int64, db *sql.DB) (*TaskSchedule, error) {
	schedules, err := GetSchedules(taskID, db)
	if err != nil {
		return nil, err
	}
	return schedules[0], nil
}

func GetSchedules(taskID int64, db *sql.DB) ([]*TaskSchedule, error) {
	query := "SELECT type, date, time, created, finished FROM task_schedule WHERE task_id=$1"

	var schedules []*TaskSchedule
	if err := sqlxt.NewScanner(db.Query(query, taskID)).Scan(&schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func GetTaskAndSchedule(taskID int64, db *sql.DB) (*Task, *TaskSchedule, error) {
	query := `SELECT t.*, s.type, s.date, s.time, s.created, s.finished
						FROM tasks AS t
						JOIN task_schedule AS s ON t.id = s.task_id
						WHERE t.id = $1`

	task := struct {
		Task     Task
		Schedule TaskSchedule
	}{}
	if err := sqlxt.NewScanner(db.Query(query, taskID)).Scan(&task); err != nil {
		return nil, nil, err
	}

	return &task.Task, &task.Schedule, task.Task.getLabels(db)
}
