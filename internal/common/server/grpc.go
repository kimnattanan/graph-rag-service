package server

import (
	"net"

	"github.com/kimnattanan/graph-rag-service/internal/common/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunGRPCServer(cfg *config.CommonConfig, register func(server *grpc.Server)) {
	RunGRPCServerOnAddr(":"+cfg.Port, register)
}

func RunGRPCServerOnAddr(addr string, register func(server *grpc.Server)) {
	grpcServer := grpc.NewServer()
	register(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.WithError(err).Panic("Unable to start gRPC server")
	}

	logrus.Info("Starting gRPC server")

	err = grpcServer.Serve(listen)
	if err != nil {
		logrus.WithError(err).Panic("Unable to start gRPC server")
	}
}
