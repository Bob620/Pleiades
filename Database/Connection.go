package Database

import (
	"github.com/globalsign/mgo"
	"log"
	"github.com/globalsign/mgo/bson"
	"../WSServer/Templates"
)

type MongodbConn struct {
	Session mgo.Session
	databaseName string
	wsclientCollectionName string
	credName string
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Source string `json:"source"`
}

type Device struct {
	Id string
}

type WSClient struct {
	Login string
	Device Device
}

type OTP struct {
	Name string
	Key []byte
}

func NewMongodbSession(url string, login Login) MongodbConn {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	err = session.Login(&mgo.Credential{Username: login.Username, Password: login.Password, Source: login.Source})
	if err != nil {
		log.Fatal(err)
	}

	return MongodbConn{*session, login.Source, "wsclients", "credBackup"}
}

func (dbConn MongodbConn) PutOTP(key []byte) {
	dbConn.Session.DB(dbConn.databaseName).C(dbConn.credName).Insert(OTP{Name: "otp", Key: key})
}

func (dbConn MongodbConn) GetOTP() OTP {
	var otp OTP

	dbConn.Session.DB(dbConn.databaseName).C(dbConn.credName).Find(bson.M{"name": "otp"}).One(&otp)

	return otp
}

func (dbConn MongodbConn) InsertConnection(Conn Templates.Connection) {
	client := WSClient{Conn.Login, Device{Id: Conn.Device.Id}}

	dbConn.Session.DB(dbConn.databaseName).C(dbConn.wsclientCollectionName).Insert(client)
}

func (dbConn MongodbConn) UpdateConnection(Conn Templates.Connection) {
	client := WSClient{Conn.Login, Device{Id: Conn.Device.Id}}

	dbConn.Session.DB(dbConn.databaseName).C(dbConn.wsclientCollectionName).Update(bson.M{"device.id": client.Device.Id}, client)
}

func (dbConn MongodbConn) GetConnection(deviceId string) Templates.Connection {
	var client WSClient

	dbConn.Session.DB(dbConn.databaseName).C(dbConn.wsclientCollectionName).Find(bson.M{"device.id": deviceId}).One(&client)

	return Templates.Connection{Device: Templates.Device{Id: client.Device.Id}, Login: client.Login}
}