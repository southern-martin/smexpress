package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Decode reads JSON body into the given struct.
func Decode(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("empty request body")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

// QueryInt extracts an integer query parameter with a default value.
func QueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

// QueryString extracts a string query parameter with a default value.
func QueryString(r *http.Request, key string, defaultVal string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}
