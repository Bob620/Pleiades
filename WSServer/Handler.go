package WSServer

import (
	"net/http"
	"github.com/gorilla/websocket"
	"src/github.com/sec51/twofactor"
	"log"
	"crypto"
	"io/ioutil"
	"src/github.com/satori/go.uuid"
	"crypto/rand"
	"../Database"
	"fmt"
	"./Templates"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
	}
)

//type Message struct {
//	Input string `json:"input"`
//}



type Handler struct {
	otp twofactor.Totp
	dbConn Database.MongodbConn
}

func (handler *Handler) Setup(dbConn Database.MongodbConn) {
	handler.dbConn = dbConn
	handler.GenerateAuth()
}

func (handler *Handler) GenerateAuth() {
	dbOTP := handler.dbConn.GetOTP()
	if dbOTP.Key == nil {
		otp, err := twofactor.NewTOTP("Pleiades Client", "Pleiades", crypto.SHA1, 6)
		if err != nil {
			log.Println("Error creating new OTP:")
			log.Fatal(err)
		}
		handler.otp = *otp

		otpBytes, err := otp.ToBytes()
		if err != nil {
			log.Fatal(err)
		}

		handler.dbConn.PutOTP(otpBytes)
	} else {
		otp, err := twofactor.TOTPFromBytes(dbOTP.Key, "Pleiades")
		if err != nil {
			log.Println("Error creating OTP from backup:")
			log.Fatal(err)
		}
		handler.otp = *otp
	}

	qrBytes, err := handler.otp.QR()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("QR.png", qrBytes, 0644)
	if err != nil {
		log.Println(err)
	}
}

func (handler *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	LoginCrypt := crypto.SHA256.New()
	connection := Templates.Connection{false, Templates.Device{""}, "", conn}

	// Authenticate a new client or login an existing one
	var initConnection Templates.InitConnection
	conn.ReadJSON(&initConnection)

	if initConnection.Auth.Key != "" {
		// Authenticate the client and provide credentials for future logins
		qrBytes, err := handler.otp.QR()
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		err = ioutil.WriteFile("test.png", qrBytes, 0644)
		if err != nil {
			log.Println(err)
		}

		err = handler.otp.Validate(initConnection.Auth.Key)
		if err == nil {
			pass, err := generateLoginCredentials()
			if err != nil {
				log.Println("Unable to generate login credentials")
				log.Println(err)
				conn.Close()
				return
			}

			deviceId, err := uuid.NewV4()
			if err != nil {
				log.Println("Unable to generate device id")
				log.Println(err)
				conn.Close()
				return
			}

			connection.Device.Id = deviceId.String()

			LoginCrypt.Reset()
			LoginCrypt.Write([]byte(connection.Device.Id+fmt.Sprintf("%x", pass)))
			connection.Login = fmt.Sprintf("%x", LoginCrypt.Sum(nil))
			connection.Authed = true

			// Store the connection
			handler.dbConn.InsertConnection(connection)

			if err := conn.WriteJSON(Templates.AuthReturn{"auth", connection.Device.Id, fmt.Sprintf("%x", pass)}); err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
			conn.Close()
			return
		}

		// Login the existing client
		var initConnection Templates.InitConnection
		conn.ReadJSON(&initConnection)
		if err != nil {
			log.Println(err)
		}

		LoginCrypt.Reset()
		LoginCrypt.Write([]byte(initConnection.Login.Id+initConnection.Login.Pass))

		if connection.Login != fmt.Sprintf("%x", LoginCrypt.Sum(nil)) {
			conn.Close()
			return
		}
	} else if initConnection.Login.Id != "" {
		// Get the existing connection info
		connection = handler.dbConn.GetConnection(initConnection.Login.Id)
		connection.Conn = conn

		// Login the existing client
		LoginCrypt.Reset()
		LoginCrypt.Write([]byte(initConnection.Login.Id+initConnection.Login.Pass))

		if connection.Login != fmt.Sprintf("%x", LoginCrypt.Sum(nil)) {
			conn.Close()
			return
		}
	} else {
		conn.Close()
		return
	}

	// Logged in, continue the connection

	if err := conn.WriteJSON(Templates.GeneralResponse{Service: "login", Type: "login", Message: "true"}); err != nil {
		log.Println(err)
	}

	for {
		var request Templates.Request
		err := conn.ReadJSON(&request)
		if err != nil {
			log.Println("Unexpected closed connection")
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Println("Unexpected closed connection")
//			}
			break
		}
	}

	// Update the connection
	conn.Close()
	connection.Conn = nil
//	handler.dbConn.UpdateConnection(connection)
	log.Println("Connection closed/device logout")
}

func generateLoginCredentials() (string, error) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
