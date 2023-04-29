# Build stage
FROM golang:1.20.3 AS build

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app .

# Final stage
FROM scratch

COPY --from=build /app /app

ENTRYPOINT ["/app"]
