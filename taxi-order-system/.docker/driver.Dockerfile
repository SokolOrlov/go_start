FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/driver/app ./app
COPY --from=build:develop /app/config/driver.docker.yaml config.yaml
COPY --from=build:develop /app/migrations/driver/ ./migrations/

CMD ["/app/app", "-c", "config.yaml"]