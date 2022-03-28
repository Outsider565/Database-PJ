FROM golang:1.16

WORKDIR /usr/src/app
COPY . /usr/src/app
EXPOSE 1234

RUN go env -w GOPROXY=https://goproxy.cn,direct\
 && go mod download

CMD go run app.go