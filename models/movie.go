package models

// Movie 电影信息
type Movie struct {
	MovieID string   `json:"movieId"`
	Title   string   `json:"title"`
	Genres  string   `json:"genres"`
	ImdbID  string   `json:"imdbId,omitempty"`
	TmdbID  string   `json:"tmdbId,omitempty"`
	Ratings []Rating `json:"ratings,omitempty"`
	Tags    []Tag    `json:"tags,omitempty"`
}

// Rating 评分信息
type Rating struct {
	UserID    string  `json:"userId"`
	Rating    float64 `json:"rating"`
	Timestamp int64   `json:"timestamp"`
}

// Tag 标签信息
type Tag struct {
	UserID    string `json:"userId"`
	Tag       string `json:"tag"`
	Timestamp int64  `json:"timestamp"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Movies []Movie `json:"movies"`
	Total  int     `json:"total"`
}
