package req

import (
	res "go/order-api/pkg/res"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, request *http.Request) (*T, error) {
	body, err := Decoder[T](request.Body)
	if err != nil {
		res.JsonResponse(*w, http.StatusBadRequest, err.Error())
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		res.JsonResponse(*w, http.StatusBadRequest, err.Error())
		return nil, err
	}
	return &body, nil
}
