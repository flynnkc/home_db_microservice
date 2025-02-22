FROM golang:1.21
WORKDIR /src
COPY . /src/
ENV CGO_ENABLED=0
RUN go mod tidy && go build -o api ./src/main.go

FROM alpine:latest
COPY --from=0 /src/api /bin/api
CMD ["/bin/api"]