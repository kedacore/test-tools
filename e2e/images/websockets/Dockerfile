FROM --platform=${BUILDPLATFORM} node:22.12.0

WORKDIR /usr/src/app

ENV GATEWAY=host.docker.internal

COPY package*.json ./
COPY server.js ./
COPY client.js ./

RUN npm install

EXPOSE 8080

CMD ["node", "server.js"]