FROM golang:1.14.4-alpine3.12 AS builder
COPY . /workspace
WORKDIR /workspace
RUN go mod download
RUN cd cmd/server && go build -o app

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
WORKDIR /workspace/
COPY --from=builder /workspace/cmd/server/app .
EXPOSE 3000
CMD ["./app"]
