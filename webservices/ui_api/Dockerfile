FROM golang:1.16.5-alpine AS BUILDER

RUN apk update && apk add --no-cache git

workdir /app

RUN git clone https://github.com/kharism/microservice_simple.git

RUN mv microservice_simple/* .

run go version

RUN go mod tidy

RUN ls

RUN cd webservices/ui_api && go build -o ui_api && cp ./ui_api ../../ && cp -r config ../../


##################################################
FROM alpine:3.13

COPY --from=builder /app/ui_api /app/ui_api

COPY --from=builder /app/config /app/config

workdir /app

RUN ls && pwd

EXPOSE 8100

ENTRYPOINT [ "/app/ui_api" ]
