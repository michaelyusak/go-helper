package dto

type Response struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
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
