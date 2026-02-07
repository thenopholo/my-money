package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func parseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}

func parseDecimal(value string) (decimal.Decimal, error) {
	return decimal.NewFromString(value)
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}

	return time.Parse("2006-01-02", value)
}

func parseDateOrNow(value string) (time.Time, error) {
	if value == "" {
		return time.Now(), nil
	}

	return parseDate(value)
}
