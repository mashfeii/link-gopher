package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
)

func ResponseWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("ResponseWithJSON", slog.Any("error", err))
	}
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	codeText := strconv.Itoa(code)
	statusText := http.StatusText(code)
	message := err.Error()

	resp := scrapper_api.ApiErrorResponse{
		Code:             &codeText,
		ExceptionName:    &statusText,
		ExceptionMessage: &message,
	}

	ResponseWithJSON(w, code, resp)
}
