FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/trip/app ./app
COPY --from=build:develop /app/config/trip.docker.yaml config.yaml
COPY --from=build:develop /app/migrations/trip/ ./migrations/

CMD ["/app/app", "-c", "config.yaml"]