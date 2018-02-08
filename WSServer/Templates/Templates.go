package Templates

import "github.com/gorilla/websocket"

type InitConnection struct {
	Auth AuthInitial `json:"auth"`
	Login Login `json:"login"`
}

type AuthInitial struct {
	DeviceType string `json:"deviceType"`
	Key string `json:"key"`
}

type AuthReturn struct {
	Type string `json:"type"`
	Id string `json:"id"`
	Pass string `json:"pass"`
}

type Login struct {
	Id string `json:"id"`
	Pass string `json:"pass"`
	RequestedServices []string `json:"requestedServices"`
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
	Login string
	Conn *websocket.Conn
}