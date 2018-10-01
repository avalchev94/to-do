package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/avalchev94/sqlxt"
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
	ID     int64       `json:"id,omitempty"`
	TaskID int64       `json:"task_id,omitempty" sql:"task_id"`
	Type   RepeatType  `json:"type"`
	Days   []int64     `json:"days"`
	Hour   clock.Clock `json:"hour"`
}

func (tr *TaskRepeatance) OK() error {
	switch {
	case tr.TaskID <= 0:
		return fmt.Errorf("TaskID is invalid")
	case !tr.Type.OK():
		return fmt.Errorf("Type is invalid")
	case len(tr.Days) == 0:
		return fmt.Errorf("Days is empty")
	case tr.Hour.String() == "":
		return fmt.Errorf("Hour is empty")
	}
	return nil
}

func (tr *TaskRepeatance) add(tx *sql.Tx) error {
	if err := tr.OK(); err != nil {
		return err
	}
	row := tx.QueryRow(`INSERT INTO task_repeatance (task_id,type,days,hour)
											VALUES ($1,$2,$3,$4) RETURNING id`,
		tr.TaskID, tr.Type, pq.Array(tr.Days), tr.Hour.HhMmSs())

	return row.Scan(&tr.ID)
}

func GetTaskRepeatance(taskID int64, db *sql.DB) (*TaskRepeatance, error) {
	row, query := db.Query("SELECT * FROM task_repeatance WHERE id=$1", taskID)

	var repeatance TaskRepeatance
	if err := sqlxt.NewScanner(row, query).Scan(&repeatance); err != nil {
		return nil, err
	}
	return &repeatance, nil
}

func GetTaskAndRepeatance(taskID int64, db *sql.DB) (*Task, *TaskRepeatance, error) {
	query := `SELECT t.*, r.type, r.days, r.hour
						FROM tasks AS t
						JOIN task_repeatance AS r ON t.id=r.task_id
						WHERE t.id=$1`
	task := struct {
		Task       Task
		Repeatance TaskRepeatance
	}{}
	if err := sqlxt.NewScanner(db.Query(query, taskID)).Scan(&task); err != nil {
		return nil, nil, err
	}

	return &task.Task, &task.Repeatance, task.Task.getLabels(db)
}
