package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrade = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		//websocket 长连接
		wsConn *websocket.Conn
		err    error
		conn   *Connection
		data   []byte
	)

	//header中添加Upgrade:websocket
	if wsConn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = InitConnection(wsConn); err != nil {
		conn.Close()
		return
	}

	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			conn.Close()
			return
		}
		if err = conn.WriteMessage(data); err != nil {
			conn.Close()
			return
		}
	}

}

func Run() {
	//http标准库
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:5555", nil)
}
