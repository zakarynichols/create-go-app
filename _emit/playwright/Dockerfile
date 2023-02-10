FROM mcr.microsoft.com/playwright:v1.30.0-jammy

WORKDIR /e2e
COPY package*.json ./

RUN npm install
COPY . .

CMD [ "npm", "run", "test" ]