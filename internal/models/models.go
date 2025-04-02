package models

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

// Request описывает запрос пользователя.
// см. https://yandex.ru/dev/dialogs/alice/doc/request.html
type RequestCreateUrl struct {
	Url string `json:"url"`
}

type ResponseCreateUrl struct {
	Result string `json:"result"`
}
