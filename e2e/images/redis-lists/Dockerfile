FROM node:12.16.2
WORKDIR /app

COPY . ./

RUN npm install --production

ENV NODE_PATH=./

ENTRYPOINT ["node", "app.js" ]