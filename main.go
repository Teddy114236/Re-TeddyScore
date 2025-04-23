package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"hbase-api/api"
	"hbase-api/config"
	"hbase-api/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.GetDefaultConfig()

	// 初始化HBase客户端
	hbaseClient := db.NewHBaseClient(cfg)
	defer hbaseClient.Close()

	// 初始化API处理器
	handler := api.NewHandler(hbaseClient)

	// 设置路由
	router := api.SetupRouter(handler)

	// 添加静态文件服务
	setupStaticFileServer(router)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("启动服务器在 %s 端口\n", cfg.Server.Port)
		log.Printf("前端界面可以通过 http://localhost:%s/frontend 访问\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v\n", err)
		}
	}()

	// 等待中断信号 关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置超时上下文以在超时时强制关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭: %v\n", err)
	}

	log.Println("服务器已关闭")
}

// 设置静态文件服务
func setupStaticFileServer(router *gin.Engine) {
	// 获取当前工作目录
	workDir, _ := os.Getwd()

	// 配置静态文件服务
	frontendPath := filepath.Join(workDir, "frontend")

	// 在/frontend路径提供frontend目录的内容
	router.StaticFS("/frontend", http.Dir(frontendPath))

	// 首页重定向到frontend/index.html
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/frontend/index.html")
	})
}
