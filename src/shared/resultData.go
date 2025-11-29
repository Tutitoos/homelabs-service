package shared

import "github.com/gofiber/fiber/v3"

type IResultError struct {
	Message string `json:"message"`
}

type IResultData[T any] struct {
	Message        string         `json:"message,omitempty"`
	Data           *T             `json:"data,omitempty"`
	Errors         []IResultError `json:"errors"`
	ContainsErrors bool           `json:"hasErrors"`
}

func ResultData[T any]() *IResultData[T] {
	return &IResultData[T]{
		Errors:         []IResultError{},
		ContainsErrors: false,
	}
}

func (r *IResultData[T]) AddMessage(message string) {
	r.Message = message
}

func (r *IResultData[T]) AddData(data T) {
	r.Data = &data
}

func (r *IResultData[T]) AddError(message string) {
	r.Errors = append(r.Errors, r.CreateError(message))
	r.ContainsErrors = true
}

func (r *IResultData[T]) AddErrors(errors []IResultError) {
	r.Errors = append(r.Errors, errors...)
	r.ContainsErrors = true
}

func (r *IResultData[T]) CreateError(message string) IResultError {
	return IResultError{
		Message: message,
	}
}

func (r *IResultData[T]) HasErrors() bool {
	return r.ContainsErrors
}

func (r *IResultData[T]) Response() fiber.Map {
	response := fiber.Map{}
	if r.Message != "" {
		response["message"] = r.Message
	}

	if r.Data != nil {
		response["data"] = *r.Data
	}

	if len(r.Errors) > 0 {
		response["errors"] = r.Errors
	}

	response["hasErrors"] = r.HasErrors()

	return response
}
