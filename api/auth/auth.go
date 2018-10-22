package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// Auth is the struct returned by New function. Can not be created explicitly.
type Auth struct {
	redis  *redis.Client
	cookie string
}

// New creates new Auth struct. The input is:
// - r: redis client for keeping the sessions ids
// - cookie: the unique name of the cookie(which keeps the session id)
func New(r *redis.Client, cookie string) *Auth {
	return &Auth{r, cookie}
}

// Login tells the "auth" package that a new user has been logged.
// A new session id and cookie (keeping the id) will be created.
func (a *Auth) Login(id int64, ctx *gin.Context) error {
	defer a.redis.Save()

	uuid := uuid.New().String()
	if err := a.redis.Set(uuid, id, 24*time.Hour).Err(); err != nil {
		return err
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  a.cookie,
		Value: uuid,
	})
	return nil
}

// Logout removes the current user's session id from the redis server.
func (a *Auth) Logout(ctx *gin.Context) error {
	uuid, err := ctx.Cookie(a.cookie)
	if err != nil {
		return nil
	}

	// delete redis record
	if err := a.redis.Del(uuid).Err(); err != nil {
		return err
	}

	// delete cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    a.cookie,
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return nil
}

// Logged returns the id(database id) of the currently logged user.
func (a *Auth) Logged(ctx *gin.Context) (int64, error) {
	uuid, err := ctx.Cookie(a.cookie)
	if err != nil {
		return 0, err
	}

	return a.redis.Get(uuid).Int64()
}
