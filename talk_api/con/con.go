package con

import (
	"log"
	"net/http"
	"talk_api/constant"

	"github.com/gorilla/websocket"
)

type header struct {
	Id int
}
type Client struct {
	Header *header
	Socket *websocket.Conn
}

type Package struct {
	Header *header
	Msg    []byte
}

type ConnectionManager struct {
	websocket.Upgrader
	clients         []*Client
	incomingPackage chan *Package
	incomingClient  chan *Client
	headerGen       <-chan *header
}

func New() ConnectionManager {
	return ConnectionManager{
		websocket.Upgrader{},
		[]*Client{},
		make(chan *Package, constant.MAX_PACKAGE_QUEUE),
		make(chan *Client, constant.MAX_CLIENT_QUEUE),
		makeHeaderGen(0),
	}
}

// Need add graceful close
func (cm ConnectionManager) AddClient(w http.ResponseWriter, req *http.Request) {
	conn, _ := cm.Upgrade(w, req, nil)
	defer conn.Close()
	head := <-cm.headerGen
	client := Client{
		head,
		conn,
	}
	cm.incomingClient <- &client
	for {
		// Receive message
		_, message, _ := conn.ReadMessage()
		pack := Package{
			head,
			message,
		}
		log.Printf("Message received: %s", pack.Msg)
		cm.incomingPackage <- &pack
	}
}

func (cm ConnectionManager) Run() {
	for {
		select {
		case pack := <-cm.incomingPackage:
			broadcast(pack, cm.clients)
		case conn := <-cm.incomingClient:
			cm.clients = append(cm.clients, conn)
		}
	}
}

// Need add graceful close
func makeHeaderGen(startValue int) <-chan *header {
	out := make(chan *header)

	go func(startValue int) {
		for {
			startValue += startValue
			head := header{
				startValue,
			}
			out <- &head
		}
	}(startValue)

	return out
}

func broadcast(pack *Package, clients []*Client) {
	// propagate message to all connections
	for _, client := range clients {
		if pack.Header != client.Header {
			_ = client.Socket.WriteMessage(websocket.TextMessage, pack.Msg)

		}
	}
}
