package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"syscall"
)

func main() {
	setRouter()
	_ = setULimit()
	log.Fatal(http.ListenAndServe(":9090", nil))
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func setULimit() error{
	var rLimit syscall.Rlimit
	if error := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); error!= nil {
		return error
	}
	rLimit.Cur = rLimit.Max
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func setRouter() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}
func reader(con *websocket.Conn) {
	for {
		messageType, p, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
			_ = con.Close()
			return
		}
		log.Println(string(p))

		if err := con.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connect success ")

	reader(ws)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello Web")
}
