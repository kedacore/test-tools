const WebSocket = require('ws');
const port = 8080;

const wss = new WebSocket.Server({ port });

const clients = new Set();

wss.on('connection', (ws) => {
  clients.add(ws);
  console.log('Client connected');

  ws.on('message', (message) => {
    console.log('Received:', message.toString());

    // Broadcast to all connected clients
    for (let client of clients) {
      if (client !== ws && client.readyState === WebSocket.OPEN) {
        client.send(message);
      }
    }
  });

  ws.on('close', () => {
    clients.delete(ws);
    console.log('Client disconnected');
  });
});


process.on('SIGINT', () => {
  console.log('\nShutting down server...');
  wss.close();
  console.log('WebSocket server closed.');
  process.exit(0);
});

console.log(`WebSocket server started on ws://localhost:${port}`);