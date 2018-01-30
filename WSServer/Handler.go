package WSServer

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
	}
)

type Handler struct {

}

func (server Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(conn.ReadMessage())

	conn.Close()
}