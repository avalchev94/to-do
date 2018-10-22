package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password,omiempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func getLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func postLogin(ctx *gin.Context) {
	username, uOK := ctx.GetPostForm("username")
	password, pOK := ctx.GetPostForm("password")
	if !uOK || !pOK {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Empty username or password field.",
		})
	}

	r, err := apiPOST("/login", &user{Name: username, Password: password}, nil)
	if err != nil {
		ctx.Error(err)
		return
	}

	switch r.StatusCode {
	case http.StatusOK:
		for _, c := range r.Cookies() {
			http.SetCookie(ctx.Writer, c)
		}
		ctx.Redirect(http.StatusSeeOther, "/")
	default:
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "some error??",
		})
	}
}

func logout(ctx *gin.Context) {
	r, err := apiGET("/logout", ctx.Request.Cookies())
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", newError(
			"Couldn't logout...",
			"Couldn't logout from the website. Try again?",
		))
		ctx.Error(err)
		return
	}

	switch r.StatusCode {
	case http.StatusOK:
		for _, c := range r.Cookies() {
			http.SetCookie(ctx.Writer, c)
		}
	default:
		ctx.Error(apiError(r).Error())
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
