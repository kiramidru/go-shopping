package catalog

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service CatalogService
}

func NewHandler(service CatalogService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	catalogGroup := rg.Group("/catalog")
	{
		catalogGroup.GET("/", h.FetchBooks)
		catalogGroup.GET("/:id", h.GetBooksByID)
	}
}

func (h *Handler) FetchBooks(c *gin.Context) {
	var req VolumeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Validation failed",
		})
		return
	}

	res, err := h.service.FetchBooks(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    "UPSTREAM_ERROR",
			"message": "Failed to search books",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetBooksByID(c *gin.Context) {
	res, err := h.service.GetBooksByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "Book not found",
			})
			return
		}

		c.JSON(http.StatusBadGateway, gin.H{
			"code":    "UPSTREAM_ERROR",
			"message": "Failed to fetch book",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
