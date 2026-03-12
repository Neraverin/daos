package handlers

import (
	"io"
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/validator"
	"github.com/gin-gonic/gin"
)

func (s *Server) ValidateRole(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, "failed to read request body")
		return
	}

	result := validator.ValidateRole(string(body))
	if result.Valid {
		valid := true
		ctx.JSON(http.StatusOK, api.ValidationResult{
			Valid:  &valid,
			Errors: &[]api.ValidationError{},
		})
		return
	}

	errors := make([]api.ValidationError, 0, len(result.Errors))
	for _, e := range result.Errors {
		errors = append(errors, api.ValidationError{
			Field:   &e.Field,
			Message: &e.Message,
		})
	}

	valid := false
	ctx.JSON(http.StatusBadRequest, api.ValidationResult{
		Valid:  &valid,
		Errors: &errors,
	})
}
