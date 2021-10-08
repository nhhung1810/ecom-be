FROM golang:1.16-alpine

WORKDIR /app
COPY ./app .

RUN go get -u github.com/gin-gonic/gin
# RUN go install -v ./...

RUN go build -o main main.go
EXPOSE 5000

CMD ["/app/main"]