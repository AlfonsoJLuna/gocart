package utils

import (
    "net/http"
    "strconv"
)

type QueryParams struct {
    Limit int
    Offset int
    OrderBy string
    OrderDir string
}

func ParseQueryParams(r *http.Request) QueryParams {
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit, _ := strconv.Atoi(limitStr)
    offset, _ := strconv.Atoi(offsetStr)

    orderBy := r.URL.Query().Get("order_by")
    orderDir := r.URL.Query().Get("order_dir")

    return QueryParams{
        Limit:    limit,
        Offset:   offset,
        OrderBy:  orderBy,
        OrderDir: orderDir,
    }
}
