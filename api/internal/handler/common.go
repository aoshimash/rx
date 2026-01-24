package handler

import (
	"fmt"
)

// parseInt parses a string to int with min/max validation
func parseInt(s string, min, max int) (int, error) {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return 0, err
	}
	if val < min || val > max {
		return 0, fmt.Errorf("value out of range")
	}
	return val, nil
}
