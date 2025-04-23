// 配置API基础URL - 自动检测当前主机
const API_BASE_URL = window.location.href.includes('localhost') || window.location.href.includes('127.0.0.1')
    ? 'http://localhost:8080/api/v1'  // 本地开发环境
    : '/api/v1';  // 生产环境（相对路径）

// Vue应用实例
new Vue({
    el: '#app',
    data: {
        // 搜索相关
        searchQuery: '',
        searchResults: [],
        searchPerformed: false,
        
        // 加载状态和错误
        loading: false,
        error: null,
        
        // 电影详情
        selectedMovie: null,
        showModal: false,
        
        // 推荐电影
        recommendedMovies: []
    },
    mounted() {
        // 添加ESC键关闭模态窗口
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.showModal) {
                this.closeModal();
            }
        });
        
        // 加载推荐电影
        this.loadRecommendedMovies();
    },
    methods: {
        // 加载推荐电影
        loadRecommendedMovies() {
            this.loading = true;
            this.error = null;
            
            // 获取ID 1-10的电影作为推荐
            const movieIds = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
            const promises = movieIds.map(id => 
                axios.get(`${API_BASE_URL}/movies/${id}`)
                    .then(response => response.data)
                    .catch(error => {
                        console.warn(`获取电影ID ${id} 失败:`, error);
                        return null;
                    })
            );
            
            Promise.all(promises)
                .then(results => {
                    // 过滤掉null结果并设置推荐电影
                    this.recommendedMovies = results.filter(movie => movie !== null);
                    this.loading = false;
                    
                    // 添加延迟让电影卡片按顺序加载
                    this.$nextTick(() => {
                        const cards = document.querySelectorAll('.recommended-movies .movie-card');
                        cards.forEach((card, index) => {
                            card.style.animationDelay = `${0.1 + index * 0.05}s`;
                        });
                    });
                })
                .catch(error => {
                    console.error('获取推荐电影出错:', error);
                    this.error = '获取推荐电影时出错，请稍后重试';
                    this.loading = false;
                });
        },
        
        // 搜索电影
        searchMovies() {
            if (!this.searchQuery.trim()) {
                return;
            }
            
            this.loading = true;
            this.error = null;
            this.searchPerformed = true;
            
            axios.get(`${API_BASE_URL}/movies`, {
                params: {
                    q: this.searchQuery,
                    limit: 20
                }
            })
            .then(response => {
                this.searchResults = response.data.movies || [];
                this.loading = false;
                
                // 添加延迟让电影卡片按顺序加载
                this.$nextTick(() => {
                    const cards = document.querySelectorAll('.movie-card');
                    cards.forEach((card, index) => {
                        card.style.animationDelay = `${0.1 + index * 0.05}s`;
                    });
                });
            })
            .catch(error => {
                console.error('搜索电影出错:', error);
                this.error = '搜索电影时出错，请稍后重试';
                this.loading = false;
            });
        },
        
        // 显示电影详情
        showMovieDetails(movieId) {
            this.loading = true;
            this.error = null;
            
            axios.get(`${API_BASE_URL}/movies/${movieId}`)
            .then(response => {
                this.selectedMovie = response.data;
                this.showModal = true;
                this.loading = false;
                
                // 添加标签动画延迟
                this.$nextTick(() => {
                    if (this.selectedMovie && this.selectedMovie.tags) {
                        const tags = document.querySelectorAll('.tag');
                        tags.forEach((tag, index) => {
                            tag.style.setProperty('--i', index);
                        });
                    }
                    
                    // 设置评分条动画
                    const bars = document.querySelectorAll('.bar');
                    bars.forEach((bar) => {
                        const width = bar.style.width;
                        bar.style.width = '0';
                        bar.style.setProperty('--width', width);
                    });
                });
            })
            .catch(error => {
                console.error('获取电影详情出错:', error);
                this.error = '获取电影详情时出错，请稍后重试';
                this.loading = false;
            });
        },
        
        // 关闭详情弹窗
        closeModal() {
            this.showModal = false;
            this.selectedMovie = null;
        },
        
        // 清除错误信息
        clearError() {
            this.error = null;
        },
        
        // 计算平均评分
        calculateAvgRating(ratings) {
            if (!ratings || ratings.length === 0) {
                return '0.0';
            }
            
            const sum = ratings.reduce((total, item) => total + item.rating, 0);
            return (sum / ratings.length).toFixed(1);
        },
        
        // 计算评分分布百分比
        calculateRatingPercentage(ratings, starLevel) {
            if (!ratings || ratings.length === 0) {
                return 0;
            }
            
            // 转换为5星制
            const fiveStarRatings = ratings.map(item => {
                // 原始评分是0-5的范围
                return Math.round(item.rating);
            });
            
            const count = fiveStarRatings.filter(rating => rating === starLevel).length;
            const percentage = Math.round((count / ratings.length) * 100);
            
            // 设置评分条宽度供CSS动画使用
            this.$nextTick(() => {
                const bars = document.querySelectorAll('.rating-bar');
                if (bars && bars[starLevel-1]) {
                    const bar = bars[starLevel-1].querySelector('.bar');
                    if (bar) {
                        bar.style.width = `${percentage}%`;
                    }
                }
            });
            
            return percentage;
        }
    }
}); 