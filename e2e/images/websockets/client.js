const WebSocket = require('ws');
const { interval } = require('rxjs');
const { take } = require('rxjs/operators');

const id = process.argv[2] || Math.floor(Math.random() * 1000);
const ws = new WebSocket(`ws://${process.env['GATEWAY'] || 'localhost'}:8080`);

ws.on('open', () => {
  console.log(`Client ${id} connected`);
  ws.send(`Hello from Client ${id}`);


  interval(4000).pipe(take(15)).subscribe((i) => {
    const message = `Message ${i + 1} from Client ${id}`;
    console.log(`Client ${id} sending: ${message}`);
    ws.send(message);
  });
});

ws.on('message', (data) => {
  const timestamp = new Date().toISOString();
  console.log(`Client ${id} received: ${data}; [${timestamp}]`);
});

ws.on('close', () => {
  console.log(`Client ${id} disconnected`);
});


process.on('SIGINT', () => {
  console.log('\nShutting down server...');
  ws.close();
  console.log('WebSocket server closed.');
  process.exit(0);
});
