package main

import (
	"log"
	"fmt"
	"flag"
	"./HttpServer"
	"./WSServer"
	"net/http"
)

var (
	wsaddr = flag.String("wsaddr", ":1337", "ws port")
	httpaddr = flag.String("httpaddr", ":8010", "http port")
)

func main() {
	var (
		webserver HttpServer.Handler
		wsserver WSServer.Handler
	)
	flag.Parse()

	webserver.Get("/", func(res http.ResponseWriter, req http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("<h1>test</h1>"))
	})

	go func() {
		log.Println(fmt.Sprintf("Now Listening for http on localhost%s", *httpaddr))
		http.ListenAndServe(*httpaddr, webserver)
	}()

	log.Println(fmt.Sprintf("Now Listening for ws on localhost%s", *wsaddr))
	http.ListenAndServe(*wsaddr, wsserver)
}
