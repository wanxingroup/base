package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/example/protos"
)

func main() {

	logger := logrus.NewEntry(logrus.New())
	logger.Logger.Level = logrus.DebugLevel
	app := launcher.NewApplication(
		launcher.SetApplicationDescription(
			&launcher.ApplicationDescription{
				Usage:            "example",
				ShortDescription: "example short description",
				LongDescription:  "example a long description for the launcher",
			},
		),
		launcher.SetApplicationLogger(logger),
		launcher.SetApplicationEvents(
			launcher.NewApplicationEvents(
				launcher.SetOnInitEvent(func(app *launcher.Application) {

					ginService := app.GetWebService()
					if ginService == nil {

						logger.WithField("stage", "onInit").Error("get gin service is nil")
						return
					}

					ginService.GetEngine().GET("/ping", pong)

					rpcService := app.GetRPCService()
					if rpcService == nil {

						logger.WithField("stage", "onInit").Error("get rpc service is nil")
						return
					}

					protos.RegisterExampleControllerServer(rpcService.GetRPCConnection(), &rpcServer{})
				}),
			),
		),
	)

	app.Launch()
}

func pong(c *gin.Context) {

	response.Response(c, "PONG")
}

type rpcServer struct {
}

func (_ rpcServer) Ping(ctx context.Context, req *protos.PingRequest) (*protos.PingReply, error) {

	if req != nil {
		return &protos.PingReply{Message: req.Message}, nil
	}

	return &protos.PingReply{Message: ""}, nil
}
