package res

import (
	"net/http"
)

func HandleBody[T any](response *http.Response) (*T, error) {
	body, err := Decoder[T](response.Body)
	if err != nil {
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
