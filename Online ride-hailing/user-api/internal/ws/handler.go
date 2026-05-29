package ws

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

// ServeWebSocket 处理WebSocket连接请求
func ServeWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL参数获取用户ID
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "缺少user_id参数", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "无效的user_id", http.StatusBadRequest)
			return
		}

		if userID <= 0 {
			http.Error(w, "无效的user_id", http.StatusBadRequest)
			return
		}

		// 升级HTTP连接为WebSocket连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("WebSocket升级失败: %v", err)
			return
		}

		// 创建客户端
		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			userID: userID,
		}

		// 注册客户端
		client.hub.register <- client

		// 启动读写协程
		go client.writePump()
		go client.readPump()

		logx.Infof("WebSocket连接成功: 用户ID=%d", userID)
	}
}
