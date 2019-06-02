package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/avalchev94/sqlxt"
)

// Label struct represents the tasks' labels.
type Label string
type Labels []Label

var (
	EmptyLabelName = errors.New("db: name can't be empty")
)

// OK validates the label data fields.
func (l Label) OK() error {
	if l == "" {
		return EmptyLabelName
	}
	return nil
}

func (l Labels) add(userID int64, tx *sql.Tx) error {
	labels := []string{}
	for _, name := range l {
		labels = append(labels, fmt.Sprintf("(%d,'%s')", userID, name))
	}

	sql := fmt.Sprintf("INSERT INTO labels (user_id, name) VALUES %s ON CONFLICT DO NOTHING",
		strings.Join(labels, ","))

	_, err := tx.Query(sql)
	return err
}

// GetLabels gets all labels for created by a user
func GetLabels(userID int64, db *sql.DB) (Labels, error) {
	labels := []struct {
		UserID int64
		Name   Label
	}{}

	scanner := sqlxt.NewScanner(db.Query("SELECT name FROM labels WHERE user_id=$1", userID))
	if err := scanner.Scan(&labels); err != nil {
		return nil, err
	}

	names := make(Labels, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return names, nil
}
