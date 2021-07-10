package util

import (
	"encoding/json"
	"net/http"
)

// WriteJSONData write http response as json
func WriteJSONData(w http.ResponseWriter, data interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	messageString := ""
	if len(message) > 0 {
		messageString = message[0]
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Success": true,
		"Data":    data,
		"Message": messageString,
	})
}

// WriteJSONDataWithTotal write http response as json with total data. Usualy it is used for grid / table data
func WriteJSONDataWithTotal(w http.ResponseWriter, data interface{}, total int64, message ...string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	messageString := ""
	if len(message) > 0 {
		messageString = message[0]
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Success": true,
		"Data":    data,
		"Total":   total,
		"Message": messageString,
	})
}

// WriteJSONError write http response error as json
func WriteJSONError(w http.ResponseWriter, err error, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if len(statusCode) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode[0])
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"Success": false,
		"Data":    nil,
		"Message": err.Error(),
	})
}
