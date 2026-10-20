package domain

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func ParseInt64QueryString(s string) ([]int64, error) {
	ids := strings.Split(s, ",")
	result := make([]int64, len(ids))
	if len(ids) == 1 && strings.TrimSpace(ids[0]) == "" {
		return []int64{}, nil
	}

	for i, id := range ids {
		idInt, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
		if err != nil {
			slog.Error("Failed to parse ID", "error", err, "id", id)
			return nil, fmt.Errorf("invalid ID: %s", id)
		}
		result[i] = idInt
	}
	return result, nil
}
