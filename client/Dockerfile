FROM golang:1.13.4

RUN go get github.com/jinzhu/gorm
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/sk409/goproxy
RUN go get github.com/sk409/goconst
RUN go get github.com/google/uuid

WORKDIR /go/src/github.com/sk409/sdso

#CMD ["go", "run", "*.go"]