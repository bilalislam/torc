package response

import "github.com/labstack/echo/v4"

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// HTTPError example
type HTTPOk struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"status ok request"`
}

func Error(ctx echo.Context, status int, err error) error {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	return ctx.JSON(status, er.Message)
}

func Result(ctx echo.Context, status int, result interface{}) error {
	return ctx.JSON(status, result)
}
