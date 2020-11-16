package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

func upper(ws *websocket.Conn) {
	var err error
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			if err != errors.New("EOF") {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println(reply)
		if err = websocket.Message.Send(ws, strings.ToUpper(reply)); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	http.Handle("/upper", websocket.Handler(upper))
	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":9999", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}