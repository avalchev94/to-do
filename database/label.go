package database

import (
	"database/sql"
	"errors"

	"github.com/avalchev94/sqlxt"
)

const (
	defaultLabelColor = "939498"
)

// Label struct represents the tasks' labels.
type Label struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id,omitempty" sql:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func isHexColor(color string) bool {
	// TODO: Add better check here
	return len(color) == 6
}

var (
	EmptyLabelName    = errors.New("db: name can't be empty")
	InvalidLabelColor = errors.New("db: invalid color")
)

// OK validates the label data fields.
func (l *Label) OK() error {
	switch {
	case l.UserID <= 0:
		return InvalidUser
	case l.Name == "":
		return EmptyLabelName
	case !isHexColor(l.Color):
		return InvalidLabelColor
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
	var label Label

	scanner := sqlxt.NewScanner(db.Query("SELECT * FROM labels WHERE id=$1", id))
	if err := scanner.Scan(&label); err != nil {
		return nil, err
	}
	return &label, nil
}

// GetLabels finds all labels for created by a user
func GetLabels(userID int64, db *sql.DB) ([]Label, error) {
	labels := []Label{}

	scanner := sqlxt.NewScanner(db.Query("SELECT * FROM labels WHERE user_id=$1", userID))
	if err := scanner.Scan(&labels); err != nil {
		return nil, err
	}
	return labels, nil
}

/*func DeleteLabel(labelID int64, db *sql.DB) error {
	db.Exec("DELETE FROM labels WHERE id")
}*/
