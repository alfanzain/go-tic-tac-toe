FROM golang:1.24.3-alpine3.21 AS build
WORKDIR /go/src/github.com/alfanzain/go-tic-tac-toe

COPY * ./
COPY ./scenes ./scenes
RUN go mod download -x
RUN go build -o game

FROM alpine:3.21
RUN apk add ca-certificates tzdata

COPY --from=build /go/src/github.com/alfanzain/go-tic-tac-toe/game ./game
COPY --from=build /go/src/github.com/alfanzain/go-tic-tac-toe/scenes ./scenes

ENTRYPOINT [ "./game" ]