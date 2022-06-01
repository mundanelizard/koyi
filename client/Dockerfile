FROM golang:1.18-alpine

WORKDIR /app

COPY ../go.mod ./
COPY ../go.sum ./

RUN go mod download

COPY *.go  ./

RUN go build -o ./client

EXPOSE 8080


FROM node:18.0-alpine

WORKDIR /app

COPY package.json ./
COPY package-lock.json ./

RUN npm ci --silent

COPY . ./

RUN npm run build

CMD [ "/client" ]
