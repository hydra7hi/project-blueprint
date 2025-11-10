# Build stage (Builds ALL services once)
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy all code
COPY . ./

# Build dependencies (runs once for all services)
RUN apk add --no-cache protoc make
RUN make clean
RUN make all
# *** 1. BUILD BOTH SERVICES HERE ***
# Build user-service
# RUN go build -o /app/grpc-user-service ./user/main.go
# Build operation-service
# RUN go build -o /app/grpc-operation-service ./operation/main.go

# Run stage for user-service (Target: user-service)
FROM alpine:latest AS user-service 
ENV SERVICE=user
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the fixed binary from the builder stage
COPY --from=builder /app/user/grpc-user-service .
EXPOSE 50051
CMD ["sh", "-c", "./grpc-user-service"]

# Run stage for operation-service (Target: operation-service)
FROM alpine:latest AS operation-service 
ENV SERVICE=operation
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the fixed binary from the builder stage
COPY --from=builder /app/operation/grpc-operation-service .
EXPOSE 50051 
CMD ["sh", "-c", "./grpc-operation-service"]