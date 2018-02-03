package main

import (
	"log"
	"fmt"
	"flag"
	"./HttpServer"
	"./WSServer"
	"net/http"
	"github.com/aymerick/raymond"
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

	webserver.AddError("404", func(res http.ResponseWriter, req http.Request) {
		ctx := map[string]string {
			"error": "404",
		}

		res.Write([]byte(raymond.MustRender(HttpServer.ReadFile("./page/error.handlebars"), ctx)))
	})

	webserver.Get("/", func(res http.ResponseWriter, req http.Request) {
		ctx := map[string]interface{} {
			"text": "test",
			"array": []map[string]interface{} {
				{"data": "Bob", "wat": "5432"},
				{"data": "Leefter", "wat": "5432"},
				{"data": "Arc", "wat": "5432"},
				{"data": "Kazuma", "wat": "5432"},
			},
		}

		res.Write([]byte(raymond.MustRender(HttpServer.ReadFile("./page/index.handlebars"), ctx)))
	})

	webserver.Routes[0].Routes = append(webserver.Routes[0].Routes, HttpServer.Route{"assets", http.MethodGet, []HttpServer.Route{}, HttpServer.ServeDir("./static/dist", "assets")})

	go func() {
		log.Println(fmt.Sprintf("Now Listening for http on localhost%s", *httpaddr))
		http.ListenAndServe(*httpaddr, webserver)
	}()

	log.Println(fmt.Sprintf("Now Listening for ws on localhost%s", *wsaddr))
	http.ListenAndServe(*wsaddr, wsserver)
}
