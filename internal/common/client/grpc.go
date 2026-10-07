package client

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClientConnection(grpcAddr string, noTLS bool) (*grpc.ClientConn, func() error, error) {
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty grpcAddr")
	}

	opts, err := grpcDialOpts(grpcAddr, noTLS)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return conn, conn.Close, nil
}

func grpcDialOpts(grpcAddr string, noTLS bool) ([]grpc.DialOption, error) {
	if noTLS {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load root CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    systemRoots,
		MinVersion: tls.VersionTLS12,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMetadataServerToken(grpcAddr)),
	}, nil
}
