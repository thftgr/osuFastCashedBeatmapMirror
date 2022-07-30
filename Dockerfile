# Run command below to build binary.
#   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main main.go

FROM scratch
WORKDIR /usr/src/app
COPY . .
CMD ["main"]