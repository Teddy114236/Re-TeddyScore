package api

import (
	"net/http"
	"strconv"

	"hbase-api/db"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	hbase *db.HBaseClient
}

// NewHandler 创建新的API处理器
func NewHandler(hbase *db.HBaseClient) *Handler {
	return &Handler{
		hbase: hbase,
	}
}

// GetMovie 获取单个电影详情
func (h *Handler) GetMovie(c *gin.Context) {
	movieID := c.Param("id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少电影ID"})
		return
	}

	movie, err := h.hbase.GetMovie(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取电影失败: " + err.Error()})
		return
	}

	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到电影"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// SearchMovies 搜索电影
func (h *Handler) SearchMovies(c *gin.Context) {
	query := c.Query("q")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	result, err := h.hbase.SearchMovies(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索电影失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
