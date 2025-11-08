package server

import (
	"grpc-services/operation/config"
	"grpc-services/operation/database"
	pb "grpc-services/operation/proto"
	userpb "grpc-services/user/proto"
)

// Server
// A struct that hold all needed values for the service during its life time.
// DB clients, and Clients of internal or external services can be added here.
// Storing Configs is usually optional.
type Server struct {
	pb.UnimplementedOperationServiceServer
	Config    *config.Config
	DB        database.DBClientInterface
	Processor *OperationProcessor
}

// NewServer
// Creates a new server instance with required dependencies.
func NewServer(cfg *config.Config, db database.DBClientInterface, userClient userpb.UserServiceClient) *Server {
	processor := NewOperationProcessor(db, userClient)

	return &Server{
		Config:    cfg,
		DB:        db,
		Processor: processor,
	}
}
