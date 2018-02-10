package main

import (
	"log"
	"fmt"
	"./HttpServer"
	"./WSServer"
	"./Database"
	"net/http"
	"github.com/aymerick/raymond"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/websocket"
	"./WSServer/Templates"
)

type Test struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Test string
	Timestamp time.Time
}

func main() {

	settings := Setup()

	var (
		webserver HttpServer.Handler
		wsserver WSServer.Handler
		dbConn Database.MongodbConn
	)

	dbConn = Database.NewMongodbSession(settings.DatabaseURL, settings.DatabaseSettings)

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

	webserver.Routes[0].Routes = append(webserver.Routes[0].Routes, HttpServer.Route{
		URL: "assets",
		Method: http.MethodGet,
		Routes: []HttpServer.Route{},
		Action: HttpServer.ServeDir("./static/dist", "assets"),
	})

	wsserver.AddService("test", func(conn *websocket.Conn,subType string, message string, variables map[string]interface{}) {
		variables["test"] = map[string]string {
			"Test": message,
		}

		switch subType {
		case "test":
			conn.WriteJSON(Templates.GeneralResponse{Service: "test", Type: "test", Message: variables["test"]})
		}
	})

	go func() {
		log.Println(fmt.Sprintf("Now Listening for http on localhost%s", ":"+settings.Httpaddr))
		http.ListenAndServe(":"+settings.Httpaddr, &webserver)
	}()

	wsserver.Setup(dbConn)

	log.Println(fmt.Sprintf("Now Listening for ws on localhost%s", ":"+settings.Wsaddr))
	http.ListenAndServe(":"+settings.Wsaddr, &wsserver)
}

type Settings struct {
	DatabaseURL string `json:"mongodbUrl"`
	DatabaseSettings Database.Login `json:"mongodbSettings"`
	Httpaddr string `json:"httpPort"`
	Wsaddr string `json:"wsPort"`
}

func Setup() Settings {
	var (
		defaultConfig Settings
		userConfig Settings
	)

	defaultFileBytes, err := ioutil.ReadFile("./config/configDefault.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(defaultFileBytes, &defaultConfig)

	fileBytes, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		ioutil.WriteFile("./config/config.json", defaultFileBytes, 0664)
	}

	json.Unmarshal(fileBytes, &userConfig)

	if userConfig.Httpaddr == "" {
		userConfig.Httpaddr = defaultConfig.Httpaddr
	}
	if userConfig.Wsaddr == "" {
		userConfig.Wsaddr = defaultConfig.Wsaddr
	}
	if userConfig.DatabaseURL == "" {
		userConfig.DatabaseURL = defaultConfig.DatabaseURL
	}

	return userConfig
}