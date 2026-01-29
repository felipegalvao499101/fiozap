package dto

import (
	"encoding/json"
	"net/http"
)

// Response resposta padrao da API
// @Description Resposta padrao da API FioZap
type Response struct {
	Code    int         `json:"code" example:"200"`
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"error message"`
}

// ActionResponse resposta generica para acoes
type ActionResponse struct {
	Details string `json:"Details" example:"Operation completed successfully"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{Code: http.StatusOK, Success: true, Data: data})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Response{Code: http.StatusCreated, Success: true, Data: data})
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Response{Code: status, Success: false, Error: msg})
}
