# TODO: Improve Dockerfile.

FROM golang:1.23.5 AS build

WORKDIR /server
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /server
COPY --from=build /server/main .
EXPOSE 3000
ENTRYPOINT ["/server/main"]