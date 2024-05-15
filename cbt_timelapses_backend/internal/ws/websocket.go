package ws

import (
	"cbt_timelapses_backend/m/v2/internal/database"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte, server *Server) // хандлер новых сообщений
	RedisDB       *redis.Client
}

func CreateServer(handleMessage func(message []byte, server *Server)) *Server {

	rdb := database.StartClient()

	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
		rdb,
	}

	http.HandleFunc("/ws", server.echo)
	http.HandleFunc("/", homePage) // delete after add frontend

	go http.ListenAndServe(":5000", nil) // Уводим http сервер в горутину

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = true // Сохраняем соединение, используя его как ключ

	server.WriteMessage([]byte("database orders data"), connection)

	defer delete(server.clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		go server.handleMessage(message, server)
	}
}

func (server *Server) WriteMessageAll(message []byte) {
	for conn := range server.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (server *Server) WriteMessage(message []byte, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, message)
}

// delete after add frontend

func homePage(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "static/index.html")
}
