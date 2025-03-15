package api

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	type Dto struct {
		Description      string `json:"description"`
		Code             int    `json:"code"`
		ExceptionMessage string `json:"exceptionMessage"`
	}

	resp := Dto{
		Description:      http.StatusText(code),
		Code:             code,
		ExceptionMessage: err.Error(),
	}

	ResponseWithJSON(w, code, resp)
}
