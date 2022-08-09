package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GenericResponse struct {
	Message string
}

func UnauthorizedA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, GenericResponse{
		Message: msg,
	})
}

func BadRequestA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, GenericResponse{
		Message: err.Error(),
	})
}

func ForbiddenA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, GenericResponse{
		Message: msg,
	})
}

func NotFoundA(c *gin.Context, entity string) {
	c.AbortWithStatusJSON(http.StatusNotFound, GenericResponse{
		Message: fmt.Sprintf("%s not found", entity),
	})
}

func ServerErrorA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, GenericResponse{
		Message: err.Error(),
	})
}
