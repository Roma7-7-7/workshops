package auth

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/wallet/api"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/wallet"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const jwtAliveDuration = 5 * time.Minute

type Middleware struct {
	repo   wallet.Repository
	secret string
}

type Claims struct {
	Timezone string
	jwt.RegisteredClaims
}

const ContextKey = "auth"

type Context struct {
	JWT *Claims
}

func (c *Context) Username() string {
	return c.JWT.Subject
}

func (c *Context) UserTimezone() string {
	return c.JWT.Timezone
}

func GetContext(c *gin.Context) *Context {
	return c.MustGet(ContextKey).(*Context)
}

func (m *Middleware) Login(c *gin.Context) {
	var req api.UserPassword
	if err := c.BindJSON(&req); err != nil {
		api.BadJSONA(c)
		return
	}
	u, err := m.repo.GetUserByName(req.Name)
	if err != nil {
		zap.L().Error("get user", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}
	if u == nil {
		api.NotFoundA(c, fmt.Sprintf("user \"%s\"", req.Name))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err == bcrypt.ErrMismatchedHashAndPassword {
		api.UnauthorizedA(c, "password does not match")
		return
	} else if err != nil {
		zap.L().Error("validate password", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	now := time.Now()
	expires := now.Add(jwtAliveDuration)
	claims := &jwt.RegisteredClaims{
		Issuer:    "wallet-app",
		Subject:   req.Name,
		ExpiresAt: jwt.NewNumericDate(expires),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(m.secret))
	if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	c.SetCookie("token", tokenS, int(jwtAliveDuration.Seconds()), "/", "wallet-app", false, false)
}

func (m *Middleware) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "wallet-app", false, false)
}

func (m *Middleware) ValidateGin(c *gin.Context) {
	tokenS, err := c.Cookie("token")
	if err == http.ErrNoCookie || (err == nil && tokenS == "") {
		api.UnauthorizedA(c, "request is not authenticated")
		return
	} else if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	cl := &Claims{}
	if _, err = jwt.ParseWithClaims(tokenS, cl, m.keyFunc); err != nil {
		api.ServerErrorA(c, err)
		return
	}
	c.Set(ContextKey, &Context{
		JWT: cl,
	})
}

func (m *Middleware) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(m.secret), nil
}

func NewMiddleware(repo wallet.Repository, secret string) *Middleware {
	return &Middleware{
		repo:   repo,
		secret: secret,
	}
}
