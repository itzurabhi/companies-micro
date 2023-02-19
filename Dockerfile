from golang:1.19.5-alpine as builder
RUN apk update --no-cache && apk add --no-cache --update go gcc g++ librdkafka-dev
WORKDIR companies

COPY . .

RUN go mod download

RUN go build -tags musl -o /bin/ ./cmd/companies

from alpine:edge

COPY --from=builder /bin/companies /companies

ENTRYPOINT ["/companies"]


