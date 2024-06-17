package util

// Convert user-friendly "page/size" pagination to DB-friendly "limit/offset".
func PageSizeToLimitOffset(page, size int) (int, int) {
	limit := size
	offset := (page - 1) * limit
	return limit, offset
}
