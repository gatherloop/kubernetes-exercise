FROM --platform=linux/amd64 golang:1.18-alpine

WORKDIR ./app

COPY . .

RUN go mod tidy
RUN go build

CMD ["./kubernetes-exercise"]