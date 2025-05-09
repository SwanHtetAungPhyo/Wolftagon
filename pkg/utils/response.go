package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code    int    `json:"status_code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// SuccessResponse creates a standardized success response
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse creates a standardized error response
func ErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	resp := Response{
		Code:    status,
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return c.Status(status).JSON(resp)
}

// ValidationErrorResponse handles validator errors specifically
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, e.Error())
		}
	} else {
		errors = append(errors, err.Error())
	}

	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Code:    fiber.StatusBadRequest,
		Message: "Validation failed",
		Error:   errors,
	})
}

// NotFoundResponse for 404 errors
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Code:    fiber.StatusNotFound,
		Message: message,
	})
}

// UnauthorizedResponse for 401 errors
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Code:    fiber.StatusUnauthorized,
		Message: message,
	})
}
