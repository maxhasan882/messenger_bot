package soc

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan Message
	register   chan Subscription
	unregister chan Subscription
}

type RawData struct {
	Data        string   `json:"data"`
	Attachments []string `json:"attachments"`
	Type        string   `json:"type"`
}

type MessageData struct {
	RoomId string `json:"room_id"`
	Data   []byte `json:"data"`
}

type Subscription struct {
	Client *Client
	roomId string
	data   chan MessageData
}

type Message struct {
	data   []byte
	roomId string
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			clientsWithRoomId := h.clients[s.roomId]
			if clientsWithRoomId == nil {
				clientsWithRoomId = make(map[*Client]bool)
				h.clients[s.roomId] = clientsWithRoomId
			}
			h.clients[s.roomId][s.Client] = true

		case s := <-h.unregister:
			clientsWithRoomId := h.clients[s.roomId]
			if clientsWithRoomId != nil {
				if _, ok := clientsWithRoomId[s.Client]; ok {
					delete(clientsWithRoomId, s.Client)
					close(s.Client.send)
					if len(clientsWithRoomId) == 0 {
						delete(h.clients, s.roomId)
					}
				}
			}

		case m := <-h.broadcast:
			clientsWithRoomId := h.clients[m.roomId]
			for c := range clientsWithRoomId {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(clientsWithRoomId, c)
					if len(clientsWithRoomId) == 0 {
						delete(h.clients, m.roomId)
					}
				}
			}
		}
	}
}
