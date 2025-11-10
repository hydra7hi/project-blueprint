## 1. Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . ./
RUN apk add --no-cache protoc make
RUN make clean
RUN make all

## 2. Run stage
# user-service
FROM alpine:latest AS user-service 
ENV SERVICE=user
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/user/grpc-user-service .
EXPOSE 50051
CMD ["sh", "-c", "./grpc-user-service"]

## operation-service
FROM alpine:latest AS operation-service 
ENV SERVICE=operation
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/operation/grpc-operation-service .
EXPOSE 50051
CMD ["sh", "-c", "./grpc-operation-service"]