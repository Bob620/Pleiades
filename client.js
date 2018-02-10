const readline = require('readline');
const WebsocketClient = require('uws');
const fs = require('fs');

const ws = new WebsocketClient('ws://localhost:1337');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

ws.on('open', () => {
    fs.readFile('cred.json', (err, data) => {
        if (err) {
            rl.question('AUTHCODE: ', (code) => {
                ws.send(JSON.stringify({
                    auth: {
                        requestedServices: ['test'],
                        deviceType: 'test',
                        key: code
                    }
                }));
            })
        } else {
            const {id, pass} = JSON.parse(data);
            ws.send(JSON.stringify({
                login: {
                    id,
                    pass
                }
            }));
        }
    });
});

ws.on('close', () => {
    console.log('closed');
});

ws.on('message', (transmission) => {
    const {service, type, id, pass, message} = JSON.parse(transmission);
    switch(type) {
        case "auth":
            fs.writeFile('cred.json', JSON.stringify({id, pass}), (err) => {
                if (err) {
                    throw err;
                }
            });
            ws.send(JSON.stringify({
                login: {
                    id,
                    pass
                }
            }));
            break;
        case "login":
            console.log("Logged in");
            break;
    }
});
