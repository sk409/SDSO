FROM golang:1.13.4

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/google/uuid
RUN go get github.com/jinzhu/gorm
RUN go get github.com/sk409/gofile
RUN go get github.com/sk409/goconst
RUN go get github.com/sk409/goproxy
RUN go get github.com/sk409/gotype
RUN go get golang.org/x/crypto/bcrypt

WORKDIR /go/src/github.com/sk409/sdso

#CMD ["go", "run", "*.go"]