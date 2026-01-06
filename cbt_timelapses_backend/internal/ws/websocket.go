package ws

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/database"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	DB            *database.DB
}

func CreateServer(handleMessage func(message []byte, server *Server)) *Server {

	db := database.StartClient()

	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
		db,
	}

	http.HandleFunc("/ws", server.echo)
	
	fileServer := http.FileServer(http.Dir(configs.SCREENSHOTS_FOLDER))
	fileServerWithCors := cors.Default().Handler(fileServer)
	http.Handle("/download/", http.StripPrefix("/download/", fileServerWithCors))


	fs := http.FileServer(http.Dir("/home/blunder/bin/cbt_timelapses_new/cbt_timelapses_frontend/dist"))
	http.Handle("/", fs)

	go http.ListenAndServe(configs.PORT_SERVER, nil) // Уводим http сервер в горутину

	log.Println("Start Server at port " + configs.PORT_SERVER)

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = true // Сохраняем соединение, используя его как ключ

	json := database.GetJSONArrayValuesFromKeyPattern(server.DB, "Order:*", true)

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
