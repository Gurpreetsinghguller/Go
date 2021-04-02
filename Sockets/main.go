package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	uid "github.com/lithammer/shortuuid"
)

type Message struct {
	Roomid  string `json:"roomid"`
	Message string `json:"message"`
}

func main() {

	server, serveError := socketio.NewServer(nil)
	if serveError != nil {
		log.Fatalln(serveError)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "createRoom", func(s socketio.Conn) string {
		u := uid.New()
		s.Join(u)
		return u
	})
	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomid string) string {

		s.Join(roomid)
		return "player joined " + s.ID()
	})
	server.OnEvent("/", "chat", func(s socketio.Conn, response Message) {

		message := response.Message
		fmt.Println("Chatting")

		server.BroadcastToRoom("", response.Roomid, "receivemsg", message)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
