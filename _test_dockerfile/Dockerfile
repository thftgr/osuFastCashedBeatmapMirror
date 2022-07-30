# Run command below to build binary.
#   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main main.go

#FROM scratch
#WORKDIR /usr/src/app
#COPY . .
#EXPOSE 80
#CMD ["main"]

FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=builder /dist/main .
ENTRYPOINT ["/main"]
