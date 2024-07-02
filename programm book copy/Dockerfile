FROM golang:1.22.3

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o books-checker cmd/main.go

CMD [ "./books-checker"]
