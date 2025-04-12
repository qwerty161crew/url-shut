package models

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

type CreateURLSRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

type CreateURLSResponse struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"short_url" validate:"required,url"`
}

type RequestCreateUrl struct {
	Url string `json:"url"`
}

type ResponseCreateUrl struct {
	Result string `json:"result"`
}
