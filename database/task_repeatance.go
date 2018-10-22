package database

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/avalchev94/sqlxt"
	"github.com/rickb777/date"
	"github.com/rickb777/date/clock"
)

type RepeatType string

const (
	Weekly  RepeatType = "weekly"
	Monthly RepeatType = "monthly"
)

func (rt RepeatType) OK() bool {
	return string(rt) == string(Weekly) || string(rt) == string(Monthly)
}

type TaskRepeatance struct {
	ID           int64       `json:"id,omitempty"`
	TaskID       int64       `json:"task_id,omitempty" sql:"task_id"`
	Type         RepeatType  `json:"type"`
	Days         []int64     `json:"days"`
	Hour         clock.Clock `json:"hour"`
	LastRepeated *date.Date  `json:"last_repeated,omitempty" sql:"last_repeated"`
}

var (
	InvalidRepeatType = errors.New("db: invalid repeat type")
	EmptyDays         = errors.New("db: days can't be empty")
	EmptyHour         = errors.New("db: hour can't be empty")
)

func (tr *TaskRepeatance) OK() error {
	switch {
	case tr.TaskID <= 0:
		return InvalidTask
	case !tr.Type.OK():
		return InvalidRepeatType
	case len(tr.Days) == 0:
		return EmptyDays
	case tr.Hour.String() == "":
		return EmptyHour
	}
	return nil
}

func (tr *TaskRepeatance) add(tx *sql.Tx) error {
	if err := tr.OK(); err != nil {
		return err
	}
	row := tx.QueryRow(`INSERT INTO task_repeatance (task_id,type,days,hour)
											VALUES ($1,$2,$3,$4) RETURNING id`,
		tr.TaskID, tr.Type, pq.Array(tr.Days), tr.Hour.HhMm())

	return row.Scan(&tr.ID)
}

func GetRepeatance(taskID int64, db *sql.DB) (*TaskRepeatance, error) {
	row, query := db.Query("SELECT * FROM task_repeatance WHERE id=$1", taskID)

	var repeatance TaskRepeatance
	if err := sqlxt.NewScanner(row, query).Scan(&repeatance); err != nil {
		return nil, err
	}
	return &repeatance, nil
}
