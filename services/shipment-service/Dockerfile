FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY services/shipment-service/ services/shipment-service/
RUN cd services/shipment-service && go build -o /bin/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/server /bin/server
EXPOSE 8089
CMD ["/bin/server"]
