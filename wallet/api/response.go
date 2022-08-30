package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GenericResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func UnauthorizedA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, GenericResponse{
		Code:    http.StatusUnauthorized,
		Message: msg,
	})
}

func BadRequestA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, GenericResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}

func BadJSONA(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, GenericResponse{
		Code:    http.StatusBadRequest,
		Message: "failed to parse request body",
	})
}

func ForbiddenA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, GenericResponse{
		Code:    http.StatusForbidden,
		Message: fmt.Sprintf("%s access denied", msg),
	})
}

func NotFoundA(c *gin.Context, entity string) {
	c.AbortWithStatusJSON(http.StatusNotFound, GenericResponse{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	})
}

func ServerErrorA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, GenericResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}
