from golang:1.19.5-alpine as builder
RUN apk update --no-cache && apk add --no-cache --update go gcc g++ librdkafka-dev
WORKDIR companies

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -tags musl -o /bin/ ./cmd/companies

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -O /bin/wait-for-it.sh && chmod +x /bin/wait-for-it.sh

from alpine:latest
RUN apk update --no-cache && apk add --no-cache bash
ENV POSTGRES_HOST=${POSTGRES_HOST}
ENV POSTGRES_PORT=${POSTGRES_PORT}
COPY --from=builder /bin/companies /companies
COPY --from=builder /bin/wait-for-it.sh /
CMD /wait-for-it.sh -t 50 -h $POSTGRES_HOST -p $POSTGRES_PORT -s -- /companies


