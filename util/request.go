package util

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// URLParam get url parameter
func URLParam(r *http.Request, key string) string {
	return strings.TrimSpace(chi.URLParam(r, key))
}

// ParsePayload parse http body to respected model
// payload := new(model.User)
// _ := util.ParsePayload(r, payload)
func ParsePayload(r *http.Request, result interface{}) error {
	return json.NewDecoder(r.Body).Decode(result)
}
