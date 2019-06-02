package database

import (
	"database/sql"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/clock"

	"github.com/avalchev94/sqlxt"
)

type RepetitiveTask struct {
	Task       Task           `json:"task"`
	Repeatance TaskRepeatance `json:"repeatance"`
}

func (t *RepetitiveTask) OK() error {
	if err := t.Task.OK(); err != nil {
		return err
	}
	if err := t.Repeatance.OK(); err != nil {
		return err
	}
	return nil
}

func (t *RepetitiveTask) Add(db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	if err := t.Task.add(transaction); err != nil {
		return transaction.Rollback()
	}

	t.Repeatance.TaskID = t.Task.ID
	if err := t.Repeatance.add(transaction); err != nil {
		return transaction.Rollback()
	}

	return transaction.Commit()
}

func (t *RepetitiveTask) Schedule(date date.Date, db *sql.DB) error {
	task := &ScheduledTask{
		Task: t.Task,
		Schedule: TaskSchedule{
			Type: DueDate,
			Date: &date,
		},
	}
	transcation, err := db.Begin()
	if err != nil {
		return err
	}
	if err := task.AddInTx(transcation); err != nil {
		transcation.Rollback()
		return err
	}

	_, err = transcation.Exec("UPDATE task_repeatance SET last_repeated=$1 WHERE id=$2",
		date, t.Task.ID)
	return err
}

func GetRepetitiveTask(taskID int64, db *sql.DB) (*RepetitiveTask, error) {
	query := `SELECT t.*, r.type, r.days, r.hour
						FROM tasks AS t
						JOIN task_repeatance AS r ON t.id=r.task_id
						WHERE t.id=$1`

	var t RepetitiveTask
	if err := sqlxt.NewScanner(db.Query(query, taskID)).Scan(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

func GetRepetitiveTasksAt(time time.Time, tasks chan<- *RepetitiveTask, db *sql.DB) error {
	clock := clock.NewAt(time)
	date := date.NewAt(time)

	query :=
		`SELECT t.*, r.type, r.days, r.hour, r.last_repeated FROM tasks AS t
		 JOIN task_repeatance AS r ON t.id = r.task_id
		 WHERE (
						(r.type = 'weekly' AND r.hour = $1 AND $2 = any(r.days)) 
						 OR
						(r.type = 'monthly' AND r.hour = $1 AND $3 = any(r.days))
					 ) AND r.last_repeated <> $4`

	rows, err := db.Query(query, clock.HhMm(), date.Weekday(), date.Day(), date.String())

	if err := sqlxt.NewScanner(rows, err).Scan(&tasks); err != nil {
		return err
	}

	return nil
}
