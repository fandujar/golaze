package golaze

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, msg string, code int) {
	JSONResponse(w, map[string]string{"error": msg}, code)
}

func JSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
