package main

import (
	"database/sql"
	"errors"
)

const (
	defaultLabelColor = "939498"
)

// Label struct represents the tasks' labels.
type Label struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id,omitempty"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func isHexColor(color string) bool {
	// TODO: Add better check here
	return len(color) == 6
}

// OK validates the label data fields.
func (l *Label) OK() error {
	if l.UserID <= 0 {
		return errors.New("invalid userID input")
	}
	if l.Name == "" {
		return errors.New("label name is empty")
	}
	if !isHexColor(l.Color) {
		return errors.New("incorrect label color")
	}
	return nil
}

func (l *Label) notExist(db *sql.DB) error {
	row := db.QueryRow("SELECT COUNT(*) FROM labels WHERE id=$1 AND name=$2", l.UserID, l.Name)
	var labels int
	if err := row.Scan(&labels); err != nil {
		return err
	}
	if labels > 0 {
		return errors.New("there is already label with that name")
	}
	return nil
}

// Add creates a label in the database.
func (l *Label) Add(db *sql.DB) error {
	if err := l.OK(); err != nil {
		return err
	}
	if err := l.notExist(db); err != nil {
		return err
	}

	_, err := db.Exec("INSERT INTO labels (user_id, name, color) VALUES($1,$2,$3)",
		l.UserID, l.Name, l.Color)

	return err
}

// GetLabel finds if there is a label with the specified id and returns it.
func GetLabel(id int64, db *sql.DB) (*Label, error) {
	row := db.QueryRow("SELECT * FROM labels WHERE id=$1", id)

	var label Label
	if err := row.Scan(&label.ID, &label.UserID, &label.Name, &label.Color); err != nil {
		return nil, err
	}
	return &label, nil
}

// GetLabels finds all labels for created by a user
func GetLabels(userID int64, db *sql.DB) ([]Label, error) {
	rows, err := db.Query("SELECT * FROM labels WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		var label Label

		if err := rows.Scan(&label.ID, &label.UserID, &label.Name, &label.Color); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	return labels, nil
}

/*func DeleteLabel(labelID int64, db *sql.DB) error {
	db.Exec("DELETE FROM labels WHERE id")
}*/
