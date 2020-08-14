package socket

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn *websocket.Conn
	// 给websocket写消息的channel
	inChan chan<- []byte
	// 读取websocket的channel
	outChan   <-chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	// closeChan 状态
	isClosed bool
}

//初始化长连接
func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan<- []byte, 1000),
		outChan:   make(<-chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	//启动读协程

	return
}

//读取websocket消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//发送消息到websocket
func (conn *Connection) WriteMessage(data []byte) error {
	select {
	case conn.inChan <- data:
	case <-conn.closeChan:
		return errors.New("connection is closed")
	}
	return nil
}

//关闭连接
func (conn *Connection) Close() {
	//线程安全的Close,可重入
	conn.wsConn.Close()

	//只执行一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			conn.Close()
			return
		}
		//如果数据量过大阻塞在这里,等待inChan有空闲的位置！
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			//closeChan关闭的时候
			conn.Close()
			return

		}
	}
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			conn.Close()
			return

		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			return
		}
	}

}
