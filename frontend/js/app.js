// 配置API基础URL
const API_BASE_URL = 'http://localhost:8080/api/v1';

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
        showModal: false
    },
    methods: {
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
            // 等动画结束后清除数据
            setTimeout(() => {
                this.selectedMovie = null;
            }, 300);
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
            return Math.round((count / ratings.length) * 100);
        }
    }
}); 