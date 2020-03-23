package stream

import (
	"errors"
	"gotify_server/model"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const writeWait = 2 * time.Second

var gerr error

type client struct {
	conn    *websocket.Conn
	userID  uint
	send    chan *model.MessageExternal
	once    once
	onClose func(*client)
}

func newClient(conn *websocket.Conn, userID uint, onClose func(*client)) *client {
	return &client{
		conn:    conn,
		send:    make(chan *model.MessageExternal, 256),
		userID:  userID,
		onClose: onClose,
	}
}

func (c *client) NotifyClose() {
	log.Println(gerr)
	c.once.Do(func() {
		c.conn.Close()
		close(c.send)
		c.onClose(c)
	})
}

func (c *client) readPump(pongWait time.Duration) {
	defer c.NotifyClose()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			printWSError("ReadError", err)
			gerr = errors.New("readerror: " + err.Error())
			return
		}
	}
}

func (c *client) writePump(pingPeriod time.Duration) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.NotifyClose()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				gerr = errors.New("not ok")
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(message); err != nil {
				gerr = errors.New("writeerror: " + err.Error())
				printWSError("WriteError", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				printWSError("PingError", err)
				gerr = errors.New("pingerror: " + err.Error())
				return
			}
		}
	}
}

func printWSError(prefix string, err error) {
	if websocket.IsCloseError(err, 1000, 1001) {
		return
	}
	log.Println("WebSocket:", prefix, err)
}
