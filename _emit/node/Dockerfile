FROM node:19

WORKDIR /client
COPY package*.json ./

RUN npm install
COPY . .

EXPOSE 7777
CMD [ "npm", "run", "serve" ]