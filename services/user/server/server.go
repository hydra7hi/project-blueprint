package server

import (
	"grpc-services/user/config"
	"grpc-services/user/database"
	pb "grpc-services/user/proto"
)

// Server
// A struct that hold all needed values for the service during its life time.
// DB clients, and Clients of internal or external services can be added here.
// Storing Configs is usually optiona.
type Server struct {
	pb.UnimplementedUserServiceServer
	Config *config.Config
	DB     database.SQLClientInterface
}

// NewServer
func NewServer(cfg *config.Config, db database.SQLClientInterface) *Server {
	return &Server{
		Config: cfg,
		DB:     db,
	}
}
