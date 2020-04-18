package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetQueryStringsForPagination(c *gin.Context) (string, int, int) {
	sort := c.Query("sort")
	limitTemp := c.Query("limit")
	pageTemp := c.Query("page")

	limit, _ := strconv.Atoi(limitTemp)
	page, _ := strconv.Atoi(pageTemp)

	if sort == "" {
		sort = "timestamp"
	}
	if limit == 0 {
		limit = 15
	}
	if page == 0 {
		page = 1
	}

	return sort, limit, page
}
