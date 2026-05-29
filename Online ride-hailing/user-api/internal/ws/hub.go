package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second
	// 读取超时时间
	pongWait = 60 * time.Second
	// 发送ping的间隔
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 512
)

// Message WebSocket消息
type Message struct {
	Type      string      `json:"type"`       // 消息类型
	Timestamp int64       `json:"timestamp"`  // 时间戳
	Data      interface{} `json:"data"`       // 消息数据
	UserID    int64       `json:"-"`          // 用户ID（内部使用）
}

// Client WebSocket客户端
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	userID int64
}

// Hub WebSocket连接管理
type Hub struct {
	// 已连接的客户端
	clients map[int64]*Client
	// 注册请求
	register chan *Client
	// 注销请求
	unregister chan *Client
	// 广播消息
	broadcast chan *Message
	// 互斥锁
	mu sync.RWMutex
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
		broadcast:  make(chan *Message),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			logx.Infof("WebSocket连接已注册: 用户ID=%d", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				logx.Infof("WebSocket连接已注销: 用户ID=%d", client.userID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			if message.UserID > 0 {
				// 发送给指定用户
				h.mu.RLock()
				client, ok := h.clients[message.UserID]
				h.mu.RUnlock()
				if ok {
					select {
					case client.send <- h.marshalMessage(message):
					default:
						close(client.send)
						h.mu.Lock()
						delete(h.clients, client.userID)
						h.mu.Unlock()
					}
				}
			} else {
				// 广播给所有用户
				h.mu.RLock()
				clients := make(map[int64]*Client, len(h.clients))
				for k, v := range h.clients {
					clients[k] = v
				}
				h.mu.RUnlock()

				for userID, client := range clients {
					select {
					case client.send <- h.marshalMessage(message):
					default:
						close(client.send)
						h.mu.Lock()
						delete(h.clients, userID)
						h.mu.Unlock()
					}
				}
			}
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, messageType string, data interface{}) {
	h.broadcast <- &Message{
		Type:      messageType,
		Timestamp: time.Now().Unix(),
		Data:      data,
		UserID:    userID,
	}
}

// Broadcast 广播消息
func (h *Hub) Broadcast(messageType string, data interface{}) {
	h.broadcast <- &Message{
		Type:      messageType,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}

// marshalMessage 序列化为JSON
func (h *Hub) marshalMessage(msg *Message) []byte {
	bytes, _ := json.Marshal(msg)
	return bytes
}

// readPump 从WebSocket连接读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("WebSocket读取错误: %v", err)
			}
			break
		}
	}
}

// writePump 向WebSocket连接写消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道已关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 检查是否有更多消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
