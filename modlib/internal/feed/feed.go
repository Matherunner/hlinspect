package feed

import (
	"hlinspect/internal/logs"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connectionPool struct {
	clients    map[string]*websocket.Conn
	clientLock sync.Mutex
}

func newConnectionPool() connectionPool {
	return connectionPool{
		clients:    map[string]*websocket.Conn{},
		clientLock: sync.Mutex{},
	}
}

func (pool *connectionPool) registerListener(key string, client *websocket.Conn) {
	pool.clientLock.Lock()
	defer pool.clientLock.Unlock()
	if cl, ok := pool.clients[key]; ok {
		cl.Close()
	}
	pool.clients[key] = client
}

func (pool *connectionPool) broadcast(message []byte) {
	pool.clientLock.Lock()
	defer pool.clientLock.Unlock()
	for key, client := range pool.clients {
		err := client.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			client.Close()
			delete(pool.clients, key)
		}
	}
}

var pool = newConnectionPool()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: return true for all now, may need to change this
		return true
	},
}

// TODO: maybe use a separate logging

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.DLLLog.Infof("web socket upgrade failed: %v", err)
		return
	}

	pool.registerListener(r.RemoteAddr, conn)
}

func Broadcast(message []byte) {
	pool.broadcast(message)
}

func Serve() {
	http.HandleFunc("/ws", webSocketHandler)
	http.ListenAndServe("localhost:32001", nil)
}
