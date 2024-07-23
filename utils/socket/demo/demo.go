package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// 升级器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理WebSocket连接
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// 读取并返回消息
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Received: %s\n", msg)
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println(err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
