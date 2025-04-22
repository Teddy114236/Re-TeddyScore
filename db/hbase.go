package db

import (
	"context"
	"strconv"
	"strings"

	"hbase-api/config"
	"hbase-api/models"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
)

// HBaseClient HBase客户端包装
type HBaseClient struct {
	client gohbase.Client
	config *config.Config
}

// NewHBaseClient 创建新的HBase客户端
func NewHBaseClient(cfg *config.Config) *HBaseClient {
	return &HBaseClient{
		client: gohbase.NewClient(cfg.HBase.ZookeeperQuorum),
		config: cfg,
	}
}

// Close 关闭HBase连接
func (h *HBaseClient) Close() {
	h.client.Close()
}

// GetMovie 通过ID获取电影信息
func (h *HBaseClient) GetMovie(ctx context.Context, movieID string) (*models.Movie, error) {
	// 创建Get请求
	getRequest, err := hrpc.NewGetStr(ctx, h.config.HBase.Table, movieID)
	if err != nil {
		return nil, err
	}

	// 执行请求
	result, err := h.client.Get(getRequest)
	if err != nil {
		return nil, err
	}

	// 如果没有结果返回nil
	if len(result.Cells) == 0 {
		return nil, nil
	}

	// 解析电影信息
	movie := &models.Movie{
		MovieID: movieID,
	}

	// 处理所有单元格
	for _, cell := range result.Cells {
		family := string(cell.Family)
		qualifier := string(cell.Qualifier)
		value := string(cell.Value)

		switch family {
		case "movie":
			if qualifier == "title" {
				movie.Title = value
			} else if qualifier == "genres" {
				movie.Genres = value
			}
		case "link":
			if qualifier == "imdbId" {
				movie.ImdbID = value
			} else if qualifier == "tmdbId" {
				movie.TmdbID = value
			}
		}
	}

	// 获取评分和标签
	movie.Ratings, _ = h.getMovieRatings(ctx, movieID)
	movie.Tags, _ = h.getMovieTags(ctx, movieID)

	return movie, nil
}

// 获取电影评分
func (h *HBaseClient) getMovieRatings(ctx context.Context, movieID string) ([]models.Rating, error) {
	// 创建Scan请求获取评分
	scanRequest, err := hrpc.NewScanRangeStr(
		ctx,
		h.config.HBase.Table,
		movieID+"_", // 起始行键
		movieID+"`", // 结束行键（使用'`'作为超出'_'的下一个字符）
		hrpc.Families(map[string][]string{"rating": nil}),
	)
	if err != nil {
		return nil, err
	}

	var ratings []models.Rating
	scanner := h.client.Scan(scanRequest)
	for {
		result, err := scanner.Next()
		if err != nil {
			break
		}
		if result == nil || len(result.Cells) == 0 {
			break
		}

		// 解析行键中的用户ID
		rowKey := string(result.Cells[0].Row)
		parts := strings.Split(rowKey, "_")
		if len(parts) != 2 {
			continue
		}
		userID := parts[1]

		var rating models.Rating
		rating.UserID = userID

		for _, cell := range result.Cells {
			qualifier := string(cell.Qualifier)
			value := string(cell.Value)

			if qualifier == "rating" {
				if ratingVal, err := strconv.ParseFloat(value, 64); err == nil {
					rating.Rating = ratingVal
				}
			} else if qualifier == "timestamp" {
				if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
					rating.Timestamp = timestamp
				}
			}
		}

		ratings = append(ratings, rating)
	}

	return ratings, nil
}

// 获取电影标签
func (h *HBaseClient) getMovieTags(ctx context.Context, movieID string) ([]models.Tag, error) {
	// 创建Scan请求获取标签
	scanRequest, err := hrpc.NewScanRangeStr(
		ctx,
		h.config.HBase.Table,
		movieID+"_", // 起始行键
		movieID+"`", // 结束行键
		hrpc.Families(map[string][]string{"tag": nil}),
	)
	if err != nil {
		return nil, err
	}

	var tags []models.Tag
	scanner := h.client.Scan(scanRequest)
	for {
		result, err := scanner.Next()
		if err != nil {
			break
		}
		if result == nil || len(result.Cells) == 0 {
			break
		}

		// 解析行键中的用户ID
		rowKey := string(result.Cells[0].Row)
		parts := strings.Split(rowKey, "_")
		if len(parts) != 2 {
			continue
		}
		userID := parts[1]

		var tag models.Tag
		tag.UserID = userID

		for _, cell := range result.Cells {
			qualifier := string(cell.Qualifier)
			value := string(cell.Value)

			if qualifier == "tag" {
				tag.Tag = value
			} else if qualifier == "timestamp" {
				if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
					tag.Timestamp = timestamp
				}
			}
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// SearchMovies 根据标题搜索电影
func (h *HBaseClient) SearchMovies(ctx context.Context, title string, limit int) (*models.SearchResult, error) {
	// 创建Scan请求
	scanRequest, err := hrpc.NewScanStr(
		ctx,
		h.config.HBase.Table,
		hrpc.Families(map[string][]string{"movie": nil}),
	)
	if err != nil {
		return nil, err
	}

	result := &models.SearchResult{
		Movies: []models.Movie{},
		Total:  0,
	}

	scanner := h.client.Scan(scanRequest)
	for {
		res, err := scanner.Next()
		if err != nil {
			break
		}
		if res == nil || len(res.Cells) == 0 {
			break
		}

		movieID := string(res.Cells[0].Row)
		var movie models.Movie
		movie.MovieID = movieID

		var hasTitle bool
		var matchesSearch bool

		for _, cell := range res.Cells {
			family := string(cell.Family)
			qualifier := string(cell.Qualifier)
			value := string(cell.Value)

			if family == "movie" {
				if qualifier == "title" {
					movie.Title = value
					hasTitle = true
					// 不区分大小写搜索
					matchesSearch = strings.Contains(
						strings.ToLower(value),
						strings.ToLower(title),
					)
				} else if qualifier == "genres" {
					movie.Genres = value
				}
			}
		}

		// 如果匹配搜索并且有标题，则添加到结果中
		if matchesSearch && hasTitle {
			result.Movies = append(result.Movies, movie)
			result.Total++
			if limit > 0 && result.Total >= limit {
				break
			}
		}
	}

	return result, nil
}
