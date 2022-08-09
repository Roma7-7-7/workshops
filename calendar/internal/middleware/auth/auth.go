package auth

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

const jwtAliveDuration = 5 * time.Minute

type Middleware struct {
	repo   calendar.Repository
	secret string
}

const сontextKey = "auth"

type Context struct {
	JWT      *jwt.RegisteredClaims
	Username string
}

func GetContext(c *gin.Context) *Context {
	return c.MustGet(сontextKey).(*Context)
}

func (m *Middleware) Login(c *gin.Context) {
	var req api.UserPassword
	if err := c.BindJSON(&req); err != nil {
		api.ServerErrorA(c, err)
		return
	}
	u, err := m.repo.GetUser(req.Username)
	if err != nil {
		log.Printf("get user: %v\n", err)
		api.ServerErrorA(c, err)
		return
	}
	if u == nil {
		api.NotFoundA(c, fmt.Sprintf("user \"%s\"", u.Name))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err == bcrypt.ErrMismatchedHashAndPassword {
		api.UnauthorizedA(c, "password does not match")
		return
	} else if err != nil {
		log.Printf("validate password: %v\n", err)
		api.ServerErrorA(c, err)
		return
	}

	now := time.Now()
	expires := now.Add(jwtAliveDuration)
	claims := &jwt.RegisteredClaims{
		Issuer:    "calendar-app",
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(expires),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(m.secret))
	if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	c.SetCookie("token", tokenS, int(jwtAliveDuration.Seconds()), "/", "calendar-app", false, false)
}

func (m *Middleware) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "calendar-app", false, false)
}

func (m *Middleware) Validate(c *gin.Context) {
	tokenS, err := c.Cookie("token")
	if err == http.ErrNoCookie || (err == nil && tokenS == "") {
		api.UnauthorizedA(c, "request is not authenticated")
		return
	} else if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	cl := &jwt.RegisteredClaims{}
	if _, err = jwt.ParseWithClaims(tokenS, cl, m.keyFunc); err != nil {
		api.ServerErrorA(c, err)
		return
	}
	c.Set(сontextKey, &Context{
		JWT:      cl,
		Username: cl.Subject,
	})
}

func (m *Middleware) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(m.secret), nil
}

func NewMiddleware(repo calendar.Repository, secret string) *Middleware {
	return &Middleware{
		repo:   repo,
		secret: secret,
	}
}
