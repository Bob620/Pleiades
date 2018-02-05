const readline = require('readline');
const WebsocketClient = require('uws');

const ws = new WebsocketClient('ws://localhost:1337');

const rl = readline.createInterface({
	input: process.stdin,
	output: process.stdout
})

ws.on('open', () => {
//  ws.send(JSON.stringify({"worldId": "toka", "storyId": "general", "event": {"type": "message"}}));
	rl.question('AUTHCODE: ', (code) => {
		ws.send(JSON.stringify({
			requestedServices: ['test'],
			deviceType: 'test',
			deviceId: '45743289645',
			key: code
		}));
	})
});

ws.on('close', () => {
	console.log('closed');
});

ws.on('message', (message) => {
  console.log(message);
  ws.send(JSON.stringify({
	  id: "test",
	  pass: "test"
  }));
});
