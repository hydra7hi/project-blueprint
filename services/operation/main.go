package main

import (
	"context"
	"log"
	"net"

	"grpc-services/operation/config"
	"grpc-services/operation/database"
	"grpc-services/operation/server"
	userCl "grpc-services/user/client"

	"google.golang.org/grpc"

	pb "grpc-services/operation/proto"
)

const (
	defaultGRPCPort = "50051"
	userServiceAddr = "user-service:50051"
)

func main() {
	ctx := context.Background()

	// Load configuration - will panic if any required env vars are missing
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	userClient, err := userCl.NewGRPCClient(ctx, userServiceAddr)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
		panic(err)
	}
	defer userClient.Close()

	// Initialize database
	db, err := database.NewPostgresClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	defer db.Close()

	// Create server instance with user client
	operationServer := server.NewServer(cfg, db, userClient)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterOperationServiceServer(grpcServer, operationServer)

	// Create tables if they don't exist
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Start background processor
	operationServer.Processor.StartBackgroundProcessor(ctx)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+defaultGRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", defaultGRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
