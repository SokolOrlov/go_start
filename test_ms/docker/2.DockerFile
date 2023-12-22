FROM golang:1.21 as builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o todos ./cmd/todos/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/todos /app/todos
COPY --from=builder /build/configs/todos.yaml /app/config.yaml
# ENTRYPOINT [ "todos" ]
EXPOSE 8080
CMD ["/app/todos", "-c", "config.yaml"]



# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /build/todos /app
# COPY --from=builder /build/configs/todos.yaml /app
# EXPOSE 8080
# CMD ["/app/app", "-c", "config.yaml"]