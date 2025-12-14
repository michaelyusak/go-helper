package dto

type Response[T any] struct {
	StatusCode int                      `json:"status_code"`
	Success    bool                     `json:"success"`
	Message    string                   `json:"message"`
	Data       T                        `json:"data,omitempty"`
	Details    []ValidationErrorDetails `json:"details,omitempty"`
}

type ErrorResponse struct {
	StatusCode int                      `json:"status_code"`
	Success    bool                     `json:"success"`
	Message    string                   `json:"message"`
	Details    []ValidationErrorDetails `json:"details,omitempty"`
}

type ValidationErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type GetHealthResponse struct {
	Health string `json:"health"`
}
