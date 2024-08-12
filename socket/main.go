package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

func generateRandomString() string {
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	runess := make([]rune, 12)
	for i := range runess {
		runess[i] = alpha[rand.Intn(len(alpha))]
	}

	return string(runess)
}

func newWsHandler(
	mutex *sync.Mutex,
	db *sql.DB,
	upgrader websocket.Upgrader,
	conns map[string]*websocket.Conn,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("estabilishing connection")

		conn, err := upgrader.Upgrade(w, r, nil)
		socketKey := r.Header.Get("Sec-Websocket-Key")

		mutex.Lock()
		conns[socketKey] = conn
		mutex.Unlock()

		if err != nil {
			log.Fatalf("error %s when upgrading connection to websocket", err)
			return
		}

		qrString, err := CreateQR(db, generateRandomString())
		if err != nil {
			log.Println(err)
		}

		qrPayload := fmt.Sprintf("%s|||%s", qrString, socketKey)

		if err := conn.WriteMessage(websocket.TextMessage, []byte(qrPayload)); err != nil {
			log.Println(err)
		}

		fmt.Println(conns)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				log.Println("closing connection")
				conn.Close()
				break
			}
		}
	}
}

func NewHandleSendMessage(message chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("message") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message <- params.Get("message")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(params.Get("message")))
	}
}

func listenToMessage(
	mutex *sync.Mutex,
	message chan string,
	db *sql.DB,
	conns map[string]*websocket.Conn,
) {
	fmt.Println("run listen...")
	for msg := range message {
		fmt.Println("ini messagenya: ", msg)
		splitted := strings.Split(msg, "|||")

		mutex.Lock()
		isExist, err := ValidateQR(db, splitted[0])
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()

		if !isExist {
			continue
		}

		spesificConn, ok := conns[splitted[1]]
		if !ok {
			continue
		}

		qrString, err := CreateQR(db, generateRandomString())
		if err != nil {
			log.Println(err)
		}

		qrPayload := fmt.Sprintf("%s|||%s", qrString, splitted[1])

		if err := spesificConn.WriteMessage(websocket.TextMessage, []byte(qrPayload)); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	message := make(chan string, 100)
	mutex := new(sync.Mutex)
	conns := make(map[string]*websocket.Conn)
	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", newWsHandler(mutex, db, upgrader, conns))
	http.HandleFunc("/send-message", NewHandleSendMessage(message))
	log.Println("Starting server...")
	log.Println("Listening to :8080")

	go listenToMessage(mutex, message, db, conns)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
