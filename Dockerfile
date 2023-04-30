# Build stage
FROM golang:1.20.3 AS build
ARG TARGETARCH=amd64
ARG TARGETOS=linux

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN echo "Building for architecture: ${TARGETARCH}, OS: ${TARGETOS}"
#RUN go env
RUN CGO_ENABLED=0; GOARCH=${TARGETARCH}; GOOS=${TARGETOS}; go build -ldflags="-w -s" -o /app .

# Final stage
FROM scratch

COPY --from=build /app /app

ENTRYPOINT ["/app"]
