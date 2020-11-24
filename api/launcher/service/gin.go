package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/middleware/log/request"
)

const ServiceNameGin = "gin"

type Gin struct {
	config     *GinConfig
	engine     *gin.Engine
	httpServer *http.Server
	logger     *logrus.Entry
	isClosing  bool
}

func NewGinService(logger *logrus.Entry, config *GinConfig) *Gin {

	g := &Gin{
		config: config,
		logger: logger,
	}

	g.engine = gin.Default()

	return g
}

func (g *Gin) OnStart() error {

	g.printConfig()
	g.initLogLevel()
	g.initRequestLogger()
	g.initHTTPServer()
	g.startHTTPServer()

	return nil
}

func (g *Gin) OnStop() error {

	g.markClosing()
	return g.closeHTTPServer()
}

func (g *Gin) GetServiceName() string {

	return ServiceNameGin
}

func (g *Gin) GetEngine() *gin.Engine {
	return g.engine
}

func (g *Gin) initLogLevel() {

	g.logger.Debug("start to init restful api service handler")
	switch g.config.GetMode() {
	case WebServiceModeDebug:
		gin.SetMode(gin.DebugMode)
	default: // Default Release Mode
		gin.SetMode(gin.ReleaseMode)
	}
	g.logger.Debug("init restful api handler service succeed")
}

func (g *Gin) initRequestLogger() {
	g.logger.Debug("start to register log middleware")
	g.engine.Use(request.ReqLoggerMiddleware())
	g.logger.Debug("register log middleware succeed")
}

func (g *Gin) printConfig() {

	g.logger.WithField("config", g.config).Info("gin config")
}

func (g *Gin) initHTTPServer() {
	g.logger.Debug("start to init restful api listener")
	// start the server，For services exposed on the public network, timeout must be set
	g.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", g.config.ListenConfig.GetIP(), g.config.ListenConfig.GetPort()),
		Handler:      g.engine,
		ReadTimeout:  g.config.ListenConfig.GetReadWriteTimeout(),
		WriteTimeout: g.config.ListenConfig.GetReadWriteTimeout(),
	}
	g.logger.Debug("init restful api listener succeed")
}

func (g *Gin) startHTTPServer() {
	g.logger.Infof("start server listening")
	go func() {
		err := g.httpServer.ListenAndServe()
		if err != nil && !g.isClosing {
			g.logger.Errorf("listen error: %v", err)
		}
	}()
}

func (g *Gin) markClosing() {
	g.isClosing = true
}

func (g *Gin) closeHTTPServer() (err error) {

	g.logger.Infof("closing http server")
	err = g.httpServer.Close()
	if err != nil {
		g.logger.Errorf("happened error at close http server: %v", err)
	}
	g.logger.Infof("http server closed")
	return
}
