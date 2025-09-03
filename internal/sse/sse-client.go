package sse

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/boyane126/go-common/config"
	"github.com/boyane126/go-common/redis"
)

// ClientManager 管理SSE客户端连接
type ClientManager struct {
	clients map[int][]chan string
	mutex   sync.RWMutex
	ctx     context.Context
	rdb     *redis.RedisClient
}

// NewClientManager 创建新的客户端管理器
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[int][]chan string),
		ctx:     context.Background(),
		rdb:     redis.Redis,
	}
}

// Initialize 初始化Redis连接并开始订阅
func (cm *ClientManager) Initialize() error {
	// 测试连接
	pong, err := cm.rdb.Client.Ping(cm.ctx).Result()
	if err != nil {
		log.Printf("Redis 连接失败: %v", err)
		return err
	}
	log.Printf("Redis 连接成功: %s", pong)

	// 开始订阅Redis消息
	go cm.subscribeRedisNotifications()

	return nil
}

// subscribeRedisNotifications 订阅Redis通知消息
func (cm *ClientManager) subscribeRedisNotifications() {
	channel := config.GetString("sse.channel")
	sub := cm.rdb.Client.Subscribe(cm.ctx, channel)
	log.Printf("已订阅 Redis 频道: %s", channel)

	ch := sub.Channel()
	log.Println("开始监听 Redis 消息...")

	for msg := range ch {
		log.Printf("收到 Redis 原始消息: %s", msg.Payload)
		var n Notification
		if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
			log.Printf("JSON 解析失败: %v, 原始数据: %s", err, msg.Payload)
			continue
		}

		log.Printf("解析后的消息: UserID=%d, Message=%s", n.UserID, n.Message)
		cm.broadcastToUser(n.UserID, n.Message)
	}
}

// broadcastToUser 向指定用户的所有客户端广播消息
func (cm *ClientManager) broadcastToUser(userID int, message string) {
	cm.mutex.RLock()
	chans, ok := cm.clients[userID]
	cm.mutex.RUnlock()

	if !ok {
		log.Printf("没有找到用户 %d 的客户端连接", userID)
		return
	}

	log.Printf("找到 %d 个客户端连接", len(chans))
	for i, c := range chans {
		log.Printf("向客户端 %d 发送消息", i)
		select {
		case c <- message:
			log.Printf("消息发送成功到客户端 %d", i)
		default:
			log.Printf("客户端 %d 消息队列已满，跳过", i)
		}
	}
}

// AddClient 添加客户端连接
func (cm *ClientManager) AddClient(userID int, ch chan string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[userID] = append(cm.clients[userID], ch)
	log.Printf("新的 SSE 连接，用户ID: %d, 当前连接数: %d", userID, len(cm.clients[userID]))
}

// RemoveClient 移除客户端连接
func (cm *ClientManager) RemoveClient(userID int, ch chan string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if chans, ok := cm.clients[userID]; ok {
		for i, c := range chans {
			if c == ch {
				cm.clients[userID] = append(cm.clients[userID][:i], cm.clients[userID][i+1:]...)
				close(ch)
				log.Printf("移除用户 %d 的客户端连接，剩余连接数: %d", userID, len(cm.clients[userID]))
				break
			}
		}
	}
}

// GetClientCount 获取指定用户的客户端连接数
func (cm *ClientManager) GetClientCount(userID int) int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	if chans, ok := cm.clients[userID]; ok {
		return len(chans)
	}
	return 0
}
