FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY services/user-service/ services/user-service/
RUN cd services/user-service && go build -o /bin/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/server /bin/server
EXPOSE 8083
CMD ["/bin/server"]
