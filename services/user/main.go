package main

import (
	"log"
	"net"
	"os"

	"grpc-services/user/config"
	"grpc-services/user/database"
	"grpc-services/user/server"

	"google.golang.org/grpc"

	pb "grpc-services/user/proto"
)

const (
	defaultGRPCPort = "50051"
)

func main() {
	// Load configuration - will panic if any required env vars are missing
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		panic(err)
	}

	// Initialize database
	db, err := database.NewPostgresClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
		panic(err)
	}

	// Create server instance
	userServer := server.NewServer(cfg, db)
	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	// Register
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+defaultGRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		panic(err)
	}

	log.Printf("gRPC server listening on port %s", defaultGRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		panic(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
