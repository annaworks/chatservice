FROM golang:1.15.0-alpine3.12 as BUILDER

ARG PORT=3000

WORKDIR /chatservice

ADD . .

RUN go build -o chatservice ./cmd/chatservice.go

FROM alpine:3.12

COPY --from=BUILDER /chatservice/chatservice /bin/.

EXPOSE $PORT

CMD ["chatservice"]
