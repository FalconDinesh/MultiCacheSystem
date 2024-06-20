package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var NotFound = errors.New("Key Does not exist")

// RespondJSON sends a JSON response with status code
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RespondError sends a JSON error response with status code
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

// RespondError sends a JSON error response with status code
// func RespondErrorWithJson(c *gin.Context, status int, err error) {
// 	c.JSON(status, gin.H{
// 		"error": err.Error(),
// 	})
// }

// LogError logs an error with a standard format
func LogError(message string, err error) {
	log.Printf("%s: %v", message, err)
}

// NotFoundHandler handles 404 Not Found errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RespondError(w, http.StatusNotFound, "Not Found")
}

// MethodNotAllowedHandler handles 405 Method Not Allowed errors
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RespondError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}
