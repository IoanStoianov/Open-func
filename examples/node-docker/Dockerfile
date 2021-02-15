FROM node:12.18.1-alpine

ENV NODE_ENV=production

WORKDIR /app

COPY . .

RUN npm install --production

CMD [ "node", "server.js" ]