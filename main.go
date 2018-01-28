
package main

import (
	"flag"
	"fmt"
	"time"
	"net/http"
	"sync/atomic"
	"github.com/gorilla/websocket"
)

type msg struct {
        str []byte
}
var loadstorevar atomic.Value
var verbose=false

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func sendServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil { fmt.Println(err) ; return }
	for {
		_, message, err := conn.ReadMessage()
		if err != nil { fmt.Println(err); return }
		v := msg{ message }
		loadstorevar.Store(v)
		if( verbose ) { fmt.Printf("mess=%s\n",message) }
	}
}

func getServer(w http.ResponseWriter, r *http.Request) {
	//var m []byte("hello")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil { fmt.Println(err) ; return }
	for {
		v := loadstorevar.Load().(msg)
		var err = conn.WriteMessage(1,v.str)
		time.Sleep(100*time.Millisecond);
		if err != nil { fmt.Println(err); return }
		if( verbose ) { fmt.Printf("sent mess\n") }
	}
}

func main() {

	// Do input args handling
	vflag:=flag.Bool("v",false,"verbose mode");
	flag.Parse();
	if( *vflag ) { verbose=true; }
	if( verbose ) { fmt.Printf("[Running in verbose mode]\n") }

	// Load a default value into the 'Store'
	// We use the store as a locking mechanism
	v := msg{ []byte("0,0,0,0,0,0,0,0,0") }
	loadstorevar.Store(v)

	// Set up each of the handlers
	http.HandleFunc("/send", sendServer)
	http.HandleFunc("/get", getServer)
	http.Handle("/files/", http.FileServer(http.Dir(".")))

	// Set the server running
	fmt.Printf("Starting\n");
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

