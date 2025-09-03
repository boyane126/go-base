package sse

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/boyane126/go-common/config"
	"github.com/gorilla/mux"
)

type Notification struct {
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

var clientManager *ClientManager

// InitSSEService 初始化SSE服务
func InitSSEService() error {
	clientManager = NewClientManager()
	return clientManager.Initialize()
}

// SSEHandler SSE 路由处理函数
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// 设置 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 处理 OPTIONS 请求
	if r.Method == "OPTIONS" {
		return
	}

	// 从查询参数获取用户ID，默认为1
	userID := 1
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = id
		}
	}

	msgChan := make(chan string, 10) // 增加缓冲区大小
	clientManager.AddClient(userID, msgChan)

	// 立即发送连接确认消息
	fmt.Fprintf(w, "data: 连接已建立\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	log.Printf("向用户 %d 发送连接确认消息", userID)

	// 使用context监听客户端断开
	ctx := r.Context()
	ticker := time.NewTicker(30 * time.Second) // 心跳
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// 客户端断开连接，清理资源
			clientManager.RemoveClient(userID, msgChan)
			return
		case <-ticker.C:
			// 发送心跳
			fmt.Fprintf(w, "data: ping\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case msg := <-msgChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

// StartSSEServer 启动SSE服务器
func StartSSEServer() {
	// 初始化SSE服务
	if err := InitSSEService(); err != nil {
		log.Fatalf("初始化SSE服务失败: %v", err)
	}

	// 创建路由
	r := mux.NewRouter()
	r.HandleFunc("/sse", SSEHandler).Methods("GET", "OPTIONS")

	// 获取配置
	host := config.GetString("sse.host")
	port := config.GetString("sse.port")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("SSE server running at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("启动SSE服务器失败: %v", err)
	}
}
