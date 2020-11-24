package service

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const ServiceNameRPC = "rpc"

type RPC struct {
	server *grpc.Server
	logger *logrus.Entry
	config *RPCConfig
}

func NewRPCService(logger *logrus.Entry, config *RPCConfig) *RPC {

	return &RPC{
		server: grpc.NewServer(),
		logger: logger,
		config: config,
	}
}

func (r *RPC) OnStart() error {

	listenAddr := fmt.Sprintf("%s:%d",
		r.config.ListenConfig.GetIP(),
		r.config.ListenConfig.GetPort())

	r.logger.WithField("listenAddr", listenAddr).Info("starting rpc service")

	listener, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return fmt.Errorf("failed to listen on %s: %s", listenAddr, err)
	}

	go func() {
		_ = r.server.Serve(listener)
	}()

	return nil
}

func (r *RPC) OnStop() error {

	r.server.GracefulStop()
	return nil
}

func (r *RPC) GetServiceName() string {

	return ServiceNameRPC
}

func (r *RPC) GetRPCConnection() *grpc.Server {

	return r.server
}
