package main
 
import (
	"golang.org/x/net/websocket"
	"net/http"
	"fmt"
)
 
func echoHandler(ws *websocket.Conn) {
	var err error
    for {
        var reply string
        if err = websocket.Message.Receive(ws, &reply); err != nil {
            fmt.Println("Can't receive")
            break
        }
        fmt.Println("Received from client: " + reply)

        msg := "Received:  " + reply
        fmt.Println("Sending to client: " + msg)

        if err = websocket.Message.Send(ws, msg); err != nil {
            fmt.Println("Can't send")
            break
        }
    }
}
 
func main() {
	fmt.Println("Open http://localhost:8090 in a browser")
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
