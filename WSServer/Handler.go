package WSServer

import (
	"net/http"
	"github.com/gorilla/websocket"
	"src/github.com/sec51/twofactor"
	"log"
	"crypto"
	"io/ioutil"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
	}
)

//type Message struct {
//	Input string `json:"input"`
//}

type AuthInitial struct {
	RequestedServices []string `json:"requestedServices"`
	DeviceType string `json:"deviceType"`
	DeviceId string `json:"deviceId"`
	Key string `json:"key"`
}

type AuthReturn struct {
	Authed bool `json:"authed"`
	Id string `json:"id"`
	Pass string `json:"pass"`
}

type Login struct {
	Id string `json:"id"`
	Pass string `json:"pass"`
}

type GeneralResponse struct {
	Service string `json:"service"`
	Type string `json:"type"`
	Message string `json:"message"`
}

type Request struct {
	Service string `json:"service"`
	Type string `json:"type"`
	Message string `json:"message"`
}

type Device struct {
	Id string
}

type Connection struct {
	Authed bool
	Device Device
	Totp []byte
	Conn *websocket.Conn
	Login string
}

type Handler struct {

}

func (handler Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	LoginCrypt := crypto.SHA256.New()

	otp, err := twofactor.NewTOTP("bruder.kraft225@gmail.com", "Bob620", crypto.SHA1, 6)
	if err != nil {
		log.Println(err)
		return
	}

	otpBytes, err := otp.ToBytes()
	if err != nil {
		log.Println(err)
		return
	}

	connection := Connection{false, Device{"test"}, otpBytes, conn, ""}

	qrBytes, err := otp.QR()
	if err != nil {
		log.Println(err)
		return
	}
	err = ioutil.WriteFile("test.png", qrBytes, 0644)
	if err != nil {
		log.Println(err)
	}

	for {
		if !connection.Authed {
			var request AuthInitial
			err := conn.ReadJSON(&request)

			if err != nil {
				log.Println(err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}
				break
			}

			err = otp.Validate(request.Key)
			if err != nil {
				log.Println(err)
				break
			}

			connection.Login = string(LoginCrypt.Sum([]byte("test"+"test")))
			connection.Authed = true

			handler.SendMessage(connection.Conn, Request{"test", "test", "test"})
		} else {
			var request Login
			err := conn.ReadJSON(&request)
			if err != nil {
				log.Println(err)
			}

			if connection.Login == string(LoginCrypt.Sum([]byte(request.Id+request.Pass))) {
				log.Println("Logged in")
			} else {
				break
			}
		}
	}

	conn.Close()
}

func (handler Handler) SendMessage(conn *websocket.Conn, request Request) {
	if err := conn.WriteJSON(request); err != nil {
		log.Println(err)
	}
}
