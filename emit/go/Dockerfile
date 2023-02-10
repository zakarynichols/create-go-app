# Go - Compile and execute go server
FROM golang:latest
WORKDIR /server
COPY . .
RUN go mod download
RUN go build -o main cmd/main.go
EXPOSE 3000
ENTRYPOINT ["/server/main"]
