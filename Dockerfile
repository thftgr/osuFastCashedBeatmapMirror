# Build stage
FROM golang:1.20.3 AS build

WORKDIR /src
#COPY go.mod go.sum ./
#RUN go get
COPY . .
RUN go env
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app .

# Final stage
FROM scratch

COPY --from=build /app /app

ENTRYPOINT ["/app"]
