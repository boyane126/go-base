# SSE (Server-Sent Events) 服务

这个项目已经成功适配了SSE功能，支持实时推送消息给客户端。

## 功能特性

- ✅ 基于Redis的消息订阅和广播
- ✅ 多用户客户端连接管理
- ✅ 自动心跳保持连接
- ✅ 优雅的客户端断开处理
- ✅ 支持CORS跨域访问
- ✅ 配置化的服务参数

## 项目结构

```
├── config/
│   └── sse.go                 # SSE服务配置
├── bootstrap/
│   └── sse.go                 # SSE服务初始化
├── internal/sse/
│   ├── service.go             # SSE HTTP服务处理
│   └── sse-client.go          # 客户端连接管理器
├── cmd/sse/
│   └── main.go                # SSE服务启动入口
└── test-sse.html              # 前端测试页面
```

## 配置说明

在环境变量中可以配置以下参数：

```bash
# SSE服务配置
SSE_HOST=0.0.0.0           # 服务监听地址，默认 0.0.0.0
SSE_PORT=8085              # 服务监听端口，默认 8085
SSE_REDIS_CHANNEL=notifications  # Redis订阅频道，默认 notifications

# Redis配置
REDIS_HOST=127.0.0.1       # Redis地址
REDIS_PORT=6379            # Redis端口
REDIS_USERNAME=            # Redis用户名（可选）
REDIS_PASSWORD=            # Redis密码（可选）
REDIS_MAIN_DB=1            # Redis数据库编号
```

## 启动服务

1. **启动SSE服务**
```bash
# 构建并运行
go build -o sse ./cmd/sse
./sse

# 或直接运行
go run ./cmd/sse
```

2. **确保Redis服务运行**
```bash
redis-server
```

## 使用方式

### 1. 客户端连接

客户端通过以下URL连接SSE服务：
```
http://localhost:8085/sse?user_id=1
```

参数说明：
- `user_id`: 用户ID，用于标识不同用户，支持多用户同时连接

### 2. 发送消息

通过Redis发布消息到指定用户：

```bash
# 使用redis-cli发送消息
redis-cli PUBLISH notifications '{"user_id":1,"message":"你好，这是一条测试消息！"}'

# 发送给不同用户
redis-cli PUBLISH notifications '{"user_id":2,"message":"发送给用户2的消息"}'
```

消息格式（JSON）：
```json
{
    "user_id": 1,
    "message": "消息内容",
    "created_at": "2024-01-01T12:00:00Z"
}
```

### 3. 前端测试

打开 `test-sse.html` 文件在浏览器中进行测试：

1. 输入用户ID
2. 点击"连接"按钮
3. 在终端中使用redis-cli发送消息
4. 观察浏览器中是否收到消息

## API接口

### GET /sse

建立SSE连接

**参数：**
- `user_id` (可选): 用户ID，默认为1

**响应：**
- Content-Type: `text/event-stream`
- 连接建立后会立即发送确认消息
- 每30秒发送一次心跳消息
- 收到Redis消息时实时推送

**示例：**
```javascript
const eventSource = new EventSource('http://localhost:8085/sse?user_id=123');

eventSource.onmessage = function(event) {
    console.log('收到消息:', event.data);
};

eventSource.onerror = function(event) {
    console.log('连接错误:', event);
};
```

## 消息流程

1. 客户端连接到 `/sse` 端点
2. 服务器将客户端注册到对应用户ID的连接池中
3. 外部系统通过Redis发布消息到 `notifications` 频道
4. SSE服务接收到Redis消息，解析用户ID
5. 将消息推送给该用户ID下的所有客户端连接
6. 客户端断开时自动清理连接资源

## 测试命令

```bash
# 1. 启动SSE服务
go run ./cmd/sse

# 2. 在另一个终端中发送测试消息
redis-cli PUBLISH notifications '{"user_id":1,"message":"测试消息1"}'
redis-cli PUBLISH notifications '{"user_id":1,"message":"测试消息2"}'
redis-cli PUBLISH notifications '{"user_id":2,"message":"发送给用户2"}'

# 3. 查看Redis频道订阅情况
redis-cli PUBSUB CHANNELS

# 4. 监控Redis消息（调试用）
redis-cli MONITOR
```

## 注意事项

1. **连接管理**: 服务会自动处理客户端断开连接，清理相关资源
2. **并发安全**: 使用读写锁保证客户端连接池的并发安全
3. **消息缓冲**: 每个客户端连接有10个消息的缓冲区，防止阻塞
4. **心跳机制**: 每30秒发送心跳，保持连接活跃
5. **CORS支持**: 支持跨域访问，方便前端集成

## 故障排查

1. **连接失败**: 检查服务是否启动，端口是否被占用
2. **收不到消息**: 检查Redis连接，确认频道名称正确
3. **消息格式错误**: 确认发送的JSON格式正确
4. **用户ID不匹配**: 检查连接时的user_id参数与消息中的user_id是否一致
