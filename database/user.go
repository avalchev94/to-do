package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/avalchev94/sqlxt"
	"golang.org/x/crypto/bcrypt"
)

// User struct describes the user information saved in the database.
type User struct {
	ID       int64     `json:"id" sql:"id"`
	Name     string    `json:"name" sql:"name"`
	Password string    `json:"password,omitempty" sql:"password"`
	Email    string    `json:"email" sql:"email"`
	Avatar   string    `json:"avatar,omitempty" sql:"avatar"`
	Created  time.Time `json:"created" sql:"created"`
}

var (
	UsernameShort = errors.New("db: short username")
	PasswordShort = errors.New("db: short password")
	EmailInvalid  = errors.New("db: email is invalid")
	EmailUsed     = errors.New("db: email is already used")
	UsernameUsed  = errors.New("db: username is already used")
)

// OK validates the user data fields.
func (u *User) OK() error {
	if len(u.Name) <= 5 {
		return UsernameShort
	}
	if len(u.Password) <= 5 {
		return PasswordShort
	}
	// TODO Add regular expression.
	if len(u.Email) <= 5 {
		return EmailInvalid
	}
	return nil
}

func (u *User) notExists(db *sql.DB) error {
	rows, err := db.Query("SELECT name, email FROM users WHERE email=$1 OR name=$2", u.Email, u.Name)

	var users []User
	if err = sqlxt.NewScanner(rows, err).Scan(&users); err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	if users[0].Name == u.Name {
		return UsernameUsed
	}
	return EmailUsed
}

// Add creates a user in the database.
func (u *User) Add(db *sql.DB) error {
	if err := u.OK(); err != nil {
		return err
	}

	if err := u.notExists(db); err != nil {
		return err
	}

	// user data is OK and does not exists in the database
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	_, err = db.Exec("INSERT INTO users (name, password, email, avatar, created) VALUES($1,$2,$3,$4,now())",
		u.Name, u.Password, u.Email, u.Avatar)

	return err
}

// VerifyLogin checks if there is a user with the specified name and password.
func VerifyLogin(name, password string, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE name=$1", name)

	var user User
	if err = sqlxt.NewScanner(rows, err).Scan(&user); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	// Don't want to expose the password, even it's encoded
	user.Password = ""
	return &user, nil
}

// GetUser finds if there is a user with the specified id and returns it.
func GetUser(id int64, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE id=$1", id)

	var user User
	if err := sqlxt.NewScanner(rows, err).Scan(&user); err != nil {
		return nil, err
	}

	// Don't want to expose the password, even it's encoded
	user.Password = ""
	return &user, nil
}
