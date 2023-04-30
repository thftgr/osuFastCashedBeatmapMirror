# Build stage
FROM golang:1.20.3 AS build
ARG TARGETARCH
ARG TARGETOS

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN echo "Building for architecture: ${TARGETARCH}, OS: ${TARGETOS}"
#RUN echo "Building for architecture: ${GOOS}, OS: ${GOARCH}"
#RUN go env GOARCH=${TARGETARCH}; GOOS=${TARGETOS}; //-ldflags="-w -s"
RUN CGO_ENABLED=0; GOARCH=${TARGETARCH}; GOOS=${TARGETOS}; go build -ldflags="-w -s" -o /app .

# Final stage
FROM scratch

COPY --from=build /app /app

ENTRYPOINT ["/app"]
