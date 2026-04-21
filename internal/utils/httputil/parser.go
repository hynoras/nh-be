package httputil

import (
	"net/http"
	"nh-be/internal/constant"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ParseStringsToUUIDs(ss []string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, s := range ss {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ParsePaginationParams(c *gin.Context) (int, int, error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}

	var err error
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid page value", err.Error())
		return 0, 0, err
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid page size value", err.Error())
		return 0, 0, err
	}
	return page, pageSize, nil
}

func ParseSortParams(c *gin.Context) (string, constant.Order, error) {
	sortBy := c.DefaultQuery("sortBy", "")
	sortOrder := c.DefaultQuery("sortOrder", "")

	if sortOrder != "" {
		sortOrder = strings.ToUpper(sortOrder)
		if sortOrder != string(constant.ASC) && sortOrder != string(constant.DESC) {
			MakeErrorResponse(c, http.StatusBadRequest, constant.ErrInvalidSortOrder,
				constant.ErrSortOrderShouldBeASCOrDESC)
			return "", "", constant.ErrSortOrderShouldBeASCOrDESC
		}
	}

	return sortBy, constant.Order(sortOrder), nil
}
