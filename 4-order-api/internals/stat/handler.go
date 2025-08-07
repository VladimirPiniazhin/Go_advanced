package stat

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/pkg/res"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.HandleFunc("GET /stat", handler.GetStat())

}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		from, err := time.Parse("2006-01-02", request.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from params", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", request.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to params", http.StatusBadRequest)
			return
		}

		by := request.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by params", http.StatusBadRequest)
			return
		}

		stats := handler.StatRepository.GetStat(by, from, to)
		fmt.Println(stats)

		res.JsonResponse(w, 200, stats)

	}

}
