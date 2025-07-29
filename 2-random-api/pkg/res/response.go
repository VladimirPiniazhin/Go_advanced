package pkg

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
