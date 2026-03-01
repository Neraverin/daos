package api

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthStatus struct {
	Status string `json:"status"`
}

func ErrorJSON(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, ErrorResponse{Error: message})
}
