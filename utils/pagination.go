package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Limit  int
	Offset int
}

func GetPaginationParams(r *http.Request, defaultLimit, defaultOffset int) PaginationParams {
	q := r.URL.Query()

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

func PaginatedResponse(data interface{}, total int64, limit, offset int) map[string]interface{} {
	return map[string]interface{}{
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"data":   data,
	}
}

func BuildFilterQuery(filters map[string]string) (string, []interface{}) {
	var queryParts []string
	var args []interface{}

	for k, v := range filters {
		queryParts = append(queryParts, fmt.Sprintf("%s LIKE ?", k))
		args = append(args, "%"+v+"%")
	}

	query := ""
	if len(queryParts) > 0 {
		query = fmt.Sprintf("%s", joinWithAND(queryParts))
	}

	return query, args
}

func joinWithAND(parts []string) string {
	return fmt.Sprintf("%s", joinStrings(parts, " AND "))
}

func joinStrings(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		result += s
		if i < len(arr)-1 {
			result += sep
		}
	}
	return result
}
