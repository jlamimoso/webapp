package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

var templates = template.Must(template.ParseFiles("echo.html"))

func inicio(w http.ResponseWriter, req *http.Request) {
	err := templates.ExecuteTemplate(w, "echo.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var quit = make(chan int)
var dado = make(chan string, 100)

func recebeUDP(end string) {
	var buf [1024]byte
	addr, _ := net.ResolveUDPAddr("udp", end)
	sock, _ := net.ListenUDP("udp", addr)
	defer sock.Close()
	m := ""
	for {
		sock.SetReadDeadline(time.Now().Add(20 * time.Second))
		rlen, err := sock.Read(buf[:])
		fmt.Printf("valor do rlen %d\n", rlen)
		if err != nil {
			fmt.Printf("Recebi um erro/timeout......%s\n", err.Error())
			m = "timeout ou fora do ar"
		} else {
			m = string(buf[0:rlen])
		}
		dado <- m
	}
}

func stop(w http.ResponseWriter, req *http.Request) {
	//	close(dado)
	http.Redirect(w, req, "/", http.StatusFound)
}

func webHandler(ws *websocket.Conn) {
	var in []byte
	//var dado make(chan string, 100)

	if err := websocket.Message.Receive(ws, &in); err != nil {
		return
	}
	fmt.Printf("Received: %s\n", string(in))

	//go recebeUDP(":9090", dado)

	i := 0
	for x := range dado {
		i++
		msg := fmt.Sprintf("recebi a msg %d %s", i, x)
		websocket.Message.Send(ws, msg)
	}
	msg := "Saindo fora!"
	websocket.Message.Send(ws, msg)
}

func main() {
	http.Handle("/", http.HandlerFunc(inicio))
	http.Handle("/stop", http.HandlerFunc(stop))
	http.HandleFunc("/echo",
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(webHandler)}
			s.ServeHTTP(w, req)
		})
	go recebeUDP(":9090")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
