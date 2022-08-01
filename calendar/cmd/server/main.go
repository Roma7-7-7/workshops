package main

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/repository/postgre"
	"github.com/Roma7-7-7/workshops/calendar/internal/server/http"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/validator"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	var dsn string
	if dsn = os.Getenv("DSN"); dsn == "" {
		dsn = "user=gouser password=gopassword dbname=gotest sslmode=disable"
	}
	repo := postgre.NewRepository(dsn)
	service := calendar.NewService(repo)
	server := http.NewServer(service, &validator.Service{})

	r := gin.Default()
	server.Register(r)
	r.Run(":5000")
}
