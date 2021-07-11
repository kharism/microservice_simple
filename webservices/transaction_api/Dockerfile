FROM golang:1.16.5-alpine AS BUILDER

RUN apk update && apk add --no-cache git

workdir /app

RUN git clone https://github.com/kharism/microservice_simple.git

RUN mv microservice_simple/* .

run go version

RUN go mod tidy

RUN ls

RUN cd webservices/transaction_api && go build -o transaction_api && cp ./transaction_api ../../ && cp -r config ../../


##################################################
FROM scratch

COPY --from=builder /app/transaction_api /app/transaction_api

COPY --from=builder /app/config /app/config

EXPOSE 8099

ENTRYPOINT [ "/app/transaction_api" ]
