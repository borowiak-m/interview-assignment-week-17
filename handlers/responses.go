package handlers

import (
	"net/http"

	"github.com/borowiak-m/interview-assignment-week-17/data"
)

// json response
func respondJSON(w http.ResponseWriter, statusCode int, apiResp *data.ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	apiResp.ToJSON(w)
}

// bad request response
func respondBadRequest(w http.ResponseWriter, message string) {
	apiResponse := &data.ApiResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}
	respondJSON(w, http.StatusBadRequest, apiResponse)
}

// internal server error response
func respondInternalServerError(w http.ResponseWriter, message string) {
	apiResponse := &data.ApiResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
	respondJSON(w, http.StatusInternalServerError, apiResponse)
}

// successful response
func respondSuccess(w http.ResponseWriter, records []*data.Record) {
	apiResponse := &data.ApiResponse{
		Code:    0,
		Message: "Success",
		Records: records,
	}
	respondJSON(w, http.StatusOK, apiResponse)
}
