package util

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/fms-ai-hackathon/model/response"
	"gorm.io/gorm"
)

var (
	pattern = "\\[(\\d+)\\]"
	regex   = regexp.MustCompile(pattern)
)

const JSON = "json"

func isEmpty(data interface{}) bool {
	if data == nil {
		return true
	}

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}

// IsDuplicateError checks if an error is due to a unique constraint violation or duplicate entry.
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error matches known duplicate-related conditions
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "already exists") ||
		strings.Contains(err.Error(), "duplicate entry") ||
		strings.Contains(err.Error(), "duplicated")
}

// WriteResponse writes a JSON response with custom HTTP code and message
func WriteResponse(ctx *gin.Context, category string, httpCode int, dataHeader interface{}, data interface{}, message string, totalPage int64, totalData int64) {
	if isEmpty(data) {
		data = []interface{}{}
	}

	if isEmpty(dataHeader) {
		dataHeader = []interface{}{}
	}

	// Check if http code is empty, then set default to 500 internal server error
	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	// Check if error message is due to a duplicate entry
	if httpCode == http.StatusInternalServerError && IsDuplicateError(errors.New(message)) {
		httpCode = http.StatusConflict // Set HTTP status code to 409 for duplicates
		message = "Conflict: Duplicate entry detected"
	}

	if category == JSON {
		ctx.JSON(
			httpCode,
			response.MainResponse{
				StatusResponse: response.StatusResponse{
					Code:    httpCode,
					Status:  http.StatusText(httpCode),
					Message: message,
				},
				TotalPage:  totalPage,
				TotalData:  totalData,
				DataHeader: dataHeader,
				Data:       data,
			},
		)
	}
}

func WriteDuplicateResponse(ctx *gin.Context, details []string) {
	WriteResponse(ctx, JSON, http.StatusConflict, []any{}, []any{}, "data duplicated", 0, 0)
}
