FROM golang:1.23 AS build
WORKDIR /app

COPY . .

RUN go build -o app .

FROM ubuntu:22.04
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/app .

RUN chmod +x ./app

EXPOSE 8080

ENTRYPOINT ["/app/app"]