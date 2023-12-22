FROM golang:alpine
WORKDIR /build
COPY . .
RUN go build -o todos cmd/todos/main.go
CMD [ "./todos" ]