package stream

import (
	"gotify_server/auth"
	"gotify_server/model"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//API provides a handler for websocket
type API struct {
	lock        sync.RWMutex
	clients     map[uint][]*client
	upgrader    *websocket.Upgrader
	pingPeriod  time.Duration
	pongTimeout time.Duration
}

//New creates a new instance of streamAPI
func New(pingPeriod, pongTimeout time.Duration) *API {
	return &API{
		clients:     make(map[uint][]*client),
		upgrader:    newUpgrader(),
		pingPeriod:  pingPeriod,
		pongTimeout: pingPeriod + pongTimeout,
	}
}

//Notify notifies the clients with given userID that a new msg received.
func (a *API) Notify(userID uint, msg *model.MessageExternal) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if clients, ok := a.clients[userID]; ok {
		for _, c := range clients {
			//may panic when sending to a closed chan
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
					}
				}()
				log.Println("sending...")
				c.send <- msg
				// default:
				// c.once.Do(func() { close(c.send) })
				// close(c.send)
			}()
		}
	}
}

//Handle handles request for upgrading
func (a *API) Handle(c *gin.Context) {
	conn, err := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	client := newClient(conn, auth.GetUserID(c), a.remove)
	a.register(client)

	go client.readPump(a.pongTimeout)
	go client.writePump(a.pingPeriod)
}

func (a *API) register(client *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.clients[client.userID] = append(a.clients[client.userID], client)
	log.Println(client.userID, ":", len(a.clients[client.userID]))
}

func (a *API) remove(c *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if clients, ok := a.clients[c.userID]; ok {
		for i, client := range clients {
			if client == c {
				a.clients[c.userID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}
}

func newUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//attention
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}
