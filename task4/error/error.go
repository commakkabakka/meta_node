package error

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 实现 Error 接口
func (e *AppError) Error() string {
	return e.Message
}

var (
	BadRequest        = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized   = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrNotFound       = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrInternalServer = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrDatabase       = &AppError{Code: http.StatusInternalServerError, Message: "database error"}
)
