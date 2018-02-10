## Documentation
As this is supposed to be used for internal networks only, the authorization is a 2auth code provided by the server and clients are stored in a mongodb

#### Methodology
Client connects, requests authorization to the server. After sending the auth code, the server will return a login id and password for further logins.

---

#### Connection in-depth
Client connects to server, server will wait for authcode to be sent
```json
{
    "auth": {
        "deviceType": "generalClient",
        "key": "AUTHCODE"
    }
}
```
After the client sends the auth code, if correct the server will send a return message with the login ID and password, otherwise close the websocket.
```json
{
    "type": "auth",
    "id": "ew9gh4e5yq3hueitgruigr",
    "pass": "ewg9hert98h54wrtht"
}
```
Once the client has this, they can then login using the credentials, and any further connections should be made using these.
```json
{
    "login": {
        "id": "ew9gh4e5yq3hueitgruigr",
        "pass": "ewg9hert98h54wrtht",
        "requestedServices": ["some service"]
    }
}
```
If the login is successful, the server will return a login event, otherwise it will close the websocket.
```json
{
    "service": "login",
    "type": "login",
    "message": "true"
}
```
At this point the client now has normal messaging protocol with the server

Example client in js: ***client.js***

---

#### Not Implemented
- Requesting services at login
- Scopes
- Normal Messaging Protocol
- Client Type Stuff

---

#### How to set up
- Install mongodb
- Create a user with read/write permission in the database you want to use
- Set up the config in config/config.json to point to and use the correct credentials

***NEVER UPLOAD YOUR CREDENTIALS TO GITHUB***