package ws

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	cors "github.com/rs/cors"
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

	// database.FlushDB(rdb)

	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
		rdb,
	}

	http.HandleFunc("/ws", server.echo)

	fileServer := http.FileServer(http.Dir(configs.SCREENSHOTS_FOLDER))
	fileServerWithCors := cors.Default().Handler(fileServer)
	http.Handle("/download/", http.StripPrefix("/download/", fileServerWithCors))

	go http.ListenAndServe(configs.PORT_SERVER, nil) // Уводим http сервер в горутину

	log.Println("Start Server at port " + configs.PORT_SERVER)

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = true // Сохраняем соединение, используя его как ключ

	json := database.GetJSONArrayValuesFromKeyPattern(server.RedisDB, "Order:*", true)

	server.WriteMessage(json, connection)

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
