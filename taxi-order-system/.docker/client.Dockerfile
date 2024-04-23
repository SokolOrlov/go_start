FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/client/app .
COPY --from=build:develop /app/config/client.docker.yaml config.yaml
COPY --from=build:develop /app/migrations/client/ ./migrations/

CMD ["/app/app", "-c", "config.yaml"]