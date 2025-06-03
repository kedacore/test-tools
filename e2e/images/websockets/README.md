# Start server for WebSocket tests

```
docker build -t ws-server .
docker run --rm -p 8080:8080 ws-server
```


# Start Clients in different terminals
```
docker run --rm ws-server node client.js 1
docker run --rm ws-server node client.js 2
docker run --rm --env GATEWAY=host.docker.internal ws-server node client.js 3
```
