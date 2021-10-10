FROM golang:1.16-alpine

ENV GOPATH=/
WORKDIR /src/ecom-be/app
COPY ./app .

RUN go get github.com/gin-gonic/gin
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/gin-contrib/cors
# RUN go install -v ./...

RUN go build -o main main.go
EXPOSE 5000

CMD ["/src/ecom-be/app/main"]